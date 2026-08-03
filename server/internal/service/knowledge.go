package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/rag"
)

const (
	metaLastIndexedAt = "knowledge.last_indexed_at"
	metaIndexStatus   = "knowledge.index_status"
	metaEmbedModel    = "knowledge.embed_model"
	metaEmbedSource   = "knowledge.embed_source"
)

type KnowledgeService struct {
	courses  *repository.CourseRepo
	knowledge *repository.KnowledgeRepo
	boot     *config.Config
	mu       sync.Mutex
	embedder rag.Embedder
	embedCfg config.EmbeddingConfig
	vectorCfg config.VectorConfig
}

func NewKnowledgeService(
	courses *repository.CourseRepo,
	knowledge *repository.KnowledgeRepo,
	cfg *config.Config,
) *KnowledgeService {
	s := &KnowledgeService{
		courses:   courses,
		knowledge: knowledge,
		boot:      cfg,
		vectorCfg: cfg.Vector,
	}
	s.reloadEmbedder(cfg.Embedding)
	return s
}

func (s *KnowledgeService) reloadEmbedder(cfg config.EmbeddingConfig) {
	if cfg.Enabled {
		s.embedder = rag.NewAPIEmbedder(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.Dimensions)
		s.embedCfg = cfg
	} else {
		s.embedder = rag.NewLocalEmbedder()
		s.embedCfg = config.EmbeddingConfig{
			Model:      "local-hash-v1",
			Dimensions: rag.LocalEmbeddingDim,
			Enabled:    false,
		}
	}
}

// UpdateEmbeddingConfig 从设置服务更新嵌入配置
func (s *KnowledgeService) UpdateEmbeddingConfig(cfg config.EmbeddingConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reloadEmbedder(cfg)
}

func (s *KnowledgeService) Status() (*model.KnowledgeStatus, error) {
	chunkCount, _ := s.knowledge.CountChunks()
	courseCount, _ := s.knowledge.CountCourses()
	chapterCount, _ := s.knowledge.CountChapters()

	lastIndexed, _ := s.knowledge.GetMeta(metaLastIndexedAt)
	indexStatus, _ := s.knowledge.GetMeta(metaIndexStatus)
	embedModel, _ := s.knowledge.GetMeta(metaEmbedModel)
	embedSource, _ := s.knowledge.GetMeta(metaEmbedSource)

	if indexStatus == "" {
		indexStatus = "idle"
	}
	if embedModel == "" {
		embedModel = s.embedder.Model()
	}
	if embedSource == "" {
		if s.embedCfg.Enabled {
			embedSource = "api"
		} else {
			embedSource = "local"
		}
	}

	var lastIndexedAt int64
	if lastIndexed != "" {
		fmt.Sscanf(lastIndexed, "%d", &lastIndexedAt)
	}

	return &model.KnowledgeStatus{
		ChunkCount:    chunkCount,
		CourseCount:   courseCount,
		ChapterCount:  chapterCount,
		EmbedModel:    embedModel,
		EmbedSource:   embedSource,
		Dimensions:    s.embedder.Dimensions(),
		LastIndexedAt: lastIndexedAt,
		IndexStatus:   indexStatus,
	}, nil
}

func (s *KnowledgeService) ListChunks(page, pageSize int) (*model.PageResult[model.KnowledgeChunk], error) {
	chunks, total, err := s.knowledge.List(page, pageSize)
	if err != nil {
		return nil, err
	}
	if chunks == nil {
		chunks = []model.KnowledgeChunk{}
	}
	for i := range chunks {
		if len(chunks[i].Text) > 200 {
			chunks[i].Text = chunks[i].Text[:200] + "…"
		}
	}
	return &model.PageResult[model.KnowledgeChunk]{
		List: chunks, Total: int(total), Page: page, PageSize: pageSize,
	}, nil
}

// RebuildIndex 全量重建知识库向量索引
func (s *KnowledgeService) RebuildIndex(ctx context.Context) (*model.KnowledgeStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	_ = s.knowledge.SetMeta(metaIndexStatus, "indexing", now)

	if err := s.knowledge.DeleteAll(); err != nil {
		_ = s.knowledge.SetMeta(metaIndexStatus, "failed", now)
		return nil, err
	}

	courses, err := s.courses.ListPublished()
	if err != nil {
		_ = s.knowledge.SetMeta(metaIndexStatus, "failed", now)
		return nil, err
	}

	embedSource := "local"
	if s.embedCfg.Enabled {
		embedSource = "api"
	}

	for _, course := range courses {
		for _, ch := range course.Chapters {
			if err := s.indexChapterLocked(ctx, course, ch, now, embedSource); err != nil {
				_ = s.knowledge.SetMeta(metaIndexStatus, "failed", now)
				return nil, err
			}
		}
	}

	_ = s.knowledge.SetMeta(metaLastIndexedAt, fmt.Sprintf("%d", now), now)
	_ = s.knowledge.SetMeta(metaIndexStatus, "ready", now)
	_ = s.knowledge.SetMeta(metaEmbedModel, s.embedder.Model(), now)
	_ = s.knowledge.SetMeta(metaEmbedSource, embedSource, now)

	return s.Status()
}

// IndexChapter 索引单个章节
func (s *KnowledgeService) IndexChapter(ctx context.Context, courseID, chapterID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	course, err := s.courses.FindByID(courseID)
	if err != nil {
		return fmt.Errorf("课程不存在")
	}
	var chapter *model.Chapter
	for i := range course.Chapters {
		if course.Chapters[i].ID == chapterID {
			chapter = &course.Chapters[i]
			break
		}
	}
	if chapter == nil {
		return fmt.Errorf("章节不存在")
	}
	if chapter.Status != "published" && course.Status != "published" {
		return s.knowledge.DeleteByChapter(courseID, chapterID)
	}

	now := time.Now().UnixMilli()
	embedSource := "local"
	if s.embedCfg.Enabled {
		embedSource = "api"
	}
	return s.indexChapterLocked(ctx, *course, *chapter, now, embedSource)
}

func (s *KnowledgeService) indexChapterLocked(ctx context.Context, course model.Course, ch model.Chapter, now int64, embedSource string) error {
	if ch.Status != "published" {
		return s.knowledge.DeleteByChapter(course.ID, ch.ID)
	}

	var chunks []model.KnowledgeChunk
	var texts []string
	for _, section := range rag.SplitByHeadings(ch.Content) {
		plain := rag.StripMarkdown(section.Body)
		if len(plain) < 30 {
			continue
		}
		hash := contentHash(course.ID, ch.ID, section.Heading, plain)
		texts = append(texts, plain)
		chunks = append(chunks, model.KnowledgeChunk{
			ID:           genID("kc"),
			CourseID:     course.ID,
			CourseTitle:  course.Title,
			ChapterID:    ch.ID,
			ChapterTitle: ch.Title,
			Heading:      section.Heading,
			Text:         plain,
			ContentHash:  hash,
			Dimensions:   s.embedder.Dimensions(),
			EmbedModel:   s.embedder.Model(),
			UpdatedAt:    now,
		})
	}
	if len(chunks) == 0 {
		return s.knowledge.DeleteByChapter(course.ID, ch.ID)
	}

	embeddings, err := s.embedder.Embed(ctx, texts)
	if err != nil {
		return err
	}
	for i := range chunks {
		chunks[i].Embedding = rag.EncodeFloat32Slice(embeddings[i])
	}
	return s.knowledge.UpsertChunks(chunks)
}

// Search 混合检索：向量相似度 + 关键词匹配
func (s *KnowledgeService) Search(ctx context.Context, query string, topK int) ([]model.KnowledgeSearchResult, error) {
	if topK <= 0 {
		topK = s.vectorCfg.TopK
	}
	if topK <= 0 {
		topK = 3
	}

	chunks, err := s.knowledge.ListAll()
	if err != nil {
		return nil, err
	}
	if len(chunks) == 0 {
		return nil, nil
	}

	queryVec, err := s.embedder.EmbedOne(ctx, query)
	if err != nil {
		return nil, err
	}

	alpha := s.vectorCfg.HybridAlpha
	if alpha <= 0 {
		alpha = 0.7
	}

	type scored struct {
		chunk        model.KnowledgeChunk
		score        float64
		vectorScore  float64
		keywordScore float64
	}
	results := make([]scored, 0, len(chunks))
	for _, ch := range chunks {
		vec := rag.DecodeFloat32Slice(ch.Embedding)
		vScore := rag.CosineSimilarity(queryVec, vec)
		kScore := rag.KeywordScore(ch.Text, ch.Heading, query)
		combined := alpha*vScore + (1-alpha)*kScore
		if combined > 0 {
			results = append(results, scored{
				chunk: ch, score: combined,
				vectorScore: vScore, keywordScore: kScore,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	if len(results) > topK {
		results = results[:topK]
	}

	out := make([]model.KnowledgeSearchResult, len(results))
	for i, r := range results {
		out[i] = model.KnowledgeSearchResult{
			Chunk: r.chunk, Score: r.score,
			VectorScore: r.vectorScore, KeywordScore: r.keywordScore,
		}
	}
	return out, nil
}

// EnsureIndexed 若索引为空则自动构建
func (s *KnowledgeService) EnsureIndexed(ctx context.Context) error {
	count, err := s.knowledge.CountChunks()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = s.RebuildIndex(ctx)
	return err
}

func contentHash(courseID, chapterID, heading, text string) string {
	h := sha256.Sum256([]byte(strings.Join([]string{courseID, chapterID, heading, text}, "|")))
	return hex.EncodeToString(h[:16])
}

// ToAISources 将检索结果转为 AI 来源引用
func ToAISources(results []model.KnowledgeSearchResult) []model.AISource {
	var sources []model.AISource
	seen := map[string]bool{}
	for _, r := range results {
		key := r.Chunk.CourseID + ":" + r.Chunk.ChapterID
		if seen[key] {
			continue
		}
		seen[key] = true
		sources = append(sources, model.AISource{
			CourseID: r.Chunk.CourseID, CourseTitle: r.Chunk.CourseTitle,
			ChapterID: r.Chunk.ChapterID, ChapterTitle: r.Chunk.ChapterTitle,
		})
	}
	return sources
}

// BuildContextFromResults 构建 RAG 上下文
func BuildContextFromResults(results []model.KnowledgeSearchResult) string {
	if len(results) == 0 {
		return "（未检索到相关知识）"
	}
	var parts []string
	for _, r := range results {
		ch := r.Chunk
		parts = append(parts, fmt.Sprintf("【%s > %s > %s】\n%s", ch.CourseTitle, ch.ChapterTitle, ch.Heading, ch.Text))
	}
	return strings.Join(parts, "\n\n")
}
