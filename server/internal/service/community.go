package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/llm"
)

// CommunityService 知识社区：问题 → 社区回答 → AI 总结 → 知识资产。
type CommunityService struct {
	repo     *repository.CommunityRepo
	users    *repository.UserRepo
	catalog  *CatalogService
	settings *SettingsService
}

func NewCommunityService(repo *repository.CommunityRepo, users *repository.UserRepo, catalog *CatalogService, settings *SettingsService) *CommunityService {
	return &CommunityService{repo: repo, users: users, catalog: catalog, settings: settings}
}

func (s *CommunityService) authors(ids []string) map[string]model.UserBrief {
	users, _ := s.users.FindByIDs(ids)
	out := map[string]model.UserBrief{}
	for _, u := range users {
		out[u.ID] = userBrief(u)
	}
	return out
}

func attachAuthor(m map[string]model.UserBrief, id string) *model.UserBrief {
	if b, ok := m[id]; ok {
		return &b
	}
	return &model.UserBrief{ID: id, Nickname: "已注销用户", Avatar: "?", AvatarColor: "#94a3b8"}
}

var validPostTypes = map[string]bool{
	model.PostQuestion: true, model.PostDiscussion: true, model.PostArticle: true, model.PostShare: true,
}

func (s *CommunityService) ListPosts(postType, keyword, authorID string, page, pageSize int) (*model.PageResult[model.CommunityPost], error) {
	page, pageSize = normalizePage(page, pageSize)
	list, total, err := s.repo.ListPosts(postType, keyword, authorID, page, pageSize)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(list))
	for _, p := range list {
		ids = append(ids, p.AuthorID)
	}
	authors := s.authors(ids)
	for i := range list {
		list[i].Author = attachAuthor(authors, list[i].AuthorID)
		list[i].Content = excerpt(list[i].Content, 140)
	}
	return &model.PageResult[model.CommunityPost]{List: nonNil(list), Total: int(total), Page: page, PageSize: pageSize}, nil
}

func excerpt(text string, n int) string {
	text = strings.Join(strings.Fields(text), " ")
	if utf8.RuneCountInString(text) <= n {
		return text
	}
	return string([]rune(text)[:n]) + "…"
}

func (s *CommunityService) GetPost(id string) (*model.PostDetail, error) {
	p, err := s.repo.FindPost(id)
	if err != nil || p.Status != "published" {
		return nil, apperr.NotFound("内容不存在")
	}
	_ = s.repo.IncrPost(id, "views", 1)
	p.Views++
	answers, _ := s.repo.ListAnswers(id)
	ids := []string{p.AuthorID}
	for _, a := range answers {
		ids = append(ids, a.AuthorID)
	}
	authors := s.authors(ids)
	p.Author = attachAuthor(authors, p.AuthorID)
	for i := range answers {
		answers[i].Author = attachAuthor(authors, answers[i].AuthorID)
	}
	d := &model.PostDetail{Post: *p, Answers: nonNil(answers), Knowledge: []model.KnowledgeBrief{}}
	if c, err := s.catalog.Load(); err == nil {
		for _, kid := range parseStrings(p.KnowledgeIDs) {
			if _, ok := c.Knowledge[kid]; ok {
				d.Knowledge = append(d.Knowledge, c.Brief(kid))
			}
		}
	}
	return d, nil
}

func (s *CommunityService) CreatePost(userID string, req model.PostRequest) (*model.CommunityPost, error) {
	req.Title, req.Content = strings.TrimSpace(req.Title), strings.TrimSpace(req.Content)
	if req.Type == "" {
		req.Type = model.PostQuestion
	}
	if !validPostTypes[req.Type] {
		return nil, apperr.BadRequest("内容类型无效")
	}
	if utf8.RuneCountInString(req.Title) < 4 || req.Content == "" {
		return nil, apperr.BadRequest("请填写标题（至少 4 个字）和正文")
	}
	now := time.Now().UnixMilli()
	p := &model.CommunityPost{
		ID: genID("post"), AuthorID: userID, Type: req.Type, Title: req.Title, Content: req.Content,
		Tags: toJSON(nonNil(req.Tags)), KnowledgeIDs: toJSON(nonNil(req.KnowledgeIDs)),
		Status: "published", CreatedAt: now, UpdatedAt: now,
	}
	return p, s.repo.SavePost(p)
}

func (s *CommunityService) Answer(userID, postID string, req model.AnswerRequest) (*model.CommunityAnswer, error) {
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return nil, apperr.BadRequest("回答内容不能为空")
	}
	if _, err := s.repo.FindPost(postID); err != nil {
		return nil, apperr.NotFound("内容不存在")
	}
	a := &model.CommunityAnswer{ID: genID("ans"), PostID: postID, AuthorID: userID, Content: req.Content, CreatedAt: time.Now().UnixMilli()}
	if err := s.repo.SaveAnswer(a); err != nil {
		return nil, err
	}
	_ = s.repo.IncrPost(postID, "answer_count", 1)
	authors := s.authors([]string{userID})
	a.Author = attachAuthor(authors, userID)
	return a, nil
}

func (s *CommunityService) LikePost(id string) error {
	if _, err := s.repo.FindPost(id); err != nil {
		return apperr.NotFound("内容不存在")
	}
	return s.repo.IncrPost(id, "likes", 1)
}

func (s *CommunityService) LikeAnswer(id string) (*model.CommunityAnswer, error) {
	a, err := s.repo.FindAnswer(id)
	if err != nil {
		return nil, apperr.NotFound("回答不存在")
	}
	a.Likes++
	return a, s.repo.SaveAnswer(a)
}

// AcceptAnswer 提问者采纳回答。
func (s *CommunityService) AcceptAnswer(userID, answerID string) error {
	a, err := s.repo.FindAnswer(answerID)
	if err != nil {
		return apperr.NotFound("回答不存在")
	}
	p, err := s.repo.FindPost(a.PostID)
	if err != nil {
		return apperr.NotFound("内容不存在")
	}
	if p.AuthorID != userID {
		return apperr.Forbidden("只有提问者可以采纳回答")
	}
	answers, _ := s.repo.ListAnswers(p.ID)
	for i := range answers {
		answers[i].Accepted = answers[i].ID == answerID
		_ = s.repo.SaveAnswer(&answers[i])
	}
	return nil
}

func (s *CommunityService) DeletePost(id string) error { return s.repo.DeletePost(id) }

// Summarize AI 总结问答讨论并沉淀为知识资产（Knowledge Resource）。
func (s *CommunityService) Summarize(ctx context.Context, userID, postID string) (*model.PostDetail, error) {
	d, err := s.GetPost(postID)
	if err != nil {
		return nil, err
	}
	if len(d.Answers) == 0 {
		return nil, apperr.BadRequest("还没有回答，暂时无法总结")
	}
	summary := s.summarize(ctx, d)
	p := d.Post
	p.AISummary = summary
	p.Author = nil
	p.UpdatedAt = time.Now().UnixMilli()

	var res *model.KnowledgeResource
	if p.ResourceID != "" {
		res, _ = s.repo.FindResource(p.ResourceID)
	}
	if res == nil {
		res = &model.KnowledgeResource{
			ID: genID("res"), Type: "article", AuthorID: p.AuthorID, SourceType: "community", SourceID: p.ID,
			Status: "published", CreatedAt: p.UpdatedAt,
		}
	}
	res.Title = "【社区沉淀】" + p.Title
	res.Summary = excerpt(summary, 120)
	res.Content = summary
	res.Tags = p.Tags
	res.KnowledgeIDs = p.KnowledgeIDs
	res.SkillIDs = toJSON([]string{})
	if err := s.repo.SaveResource(res); err != nil {
		return nil, err
	}
	p.ResourceID = res.ID
	if err := s.repo.SavePost(&p); err != nil {
		return nil, err
	}
	return s.GetPost(postID)
}

func (s *CommunityService) summarize(ctx context.Context, d *model.PostDetail) string {
	var b strings.Builder
	for i, a := range d.Answers {
		if i >= 8 {
			break
		}
		flag := ""
		if a.Accepted {
			flag = "（已采纳）"
		}
		fmt.Fprintf(&b, "回答%d%s（%d 赞）：%s\n\n", i+1, flag, a.Likes, a.Content)
	}
	if client := s.settings.LLMClient(); client != nil && client.Enabled() {
		raw, err := client.Complete(ctx, []llm.Message{
			{Role: "system", Content: "你是 Knowledge Agent，负责把社区问答整理成结构化知识资产。使用中文 Markdown，包含：核心结论、关键要点（列表）、易错点、延伸学习。不要编造回答中没有的事实。"},
			{Role: "user", Content: fmt.Sprintf("问题：%s\n\n问题描述：%s\n\n%s", d.Post.Title, d.Post.Content, b.String())},
		})
		if err == nil && strings.TrimSpace(raw) != "" {
			return strings.TrimSpace(raw)
		}
	}
	return extractiveSummary(d)
}

// extractiveSummary 未配置大模型时的抽取式总结：优先采纳回答，其次高赞回答。
func extractiveSummary(d *model.PostDetail) string {
	var b strings.Builder
	b.WriteString("## 核心结论\n\n")
	best := d.Answers[0]
	for _, a := range d.Answers {
		if a.Accepted || a.Likes > best.Likes {
			best = a
			if a.Accepted {
				break
			}
		}
	}
	b.WriteString(excerpt(best.Content, 220) + "\n\n## 社区观点\n\n")
	for i, a := range d.Answers {
		if i >= 5 {
			break
		}
		name := "社区成员"
		if a.Author != nil {
			name = a.Author.Nickname
		}
		fmt.Fprintf(&b, "- **%s**：%s\n", name, excerpt(a.Content, 80))
	}
	b.WriteString("\n> 由规则抽取生成（配置大模型后由 Knowledge Agent 生成结构化总结）")
	return b.String()
}

// ---------- 知识资产 ----------

func (s *CommunityService) ListResources(resType, keyword string, page, pageSize int) (*model.PageResult[model.KnowledgeResource], error) {
	page, pageSize = normalizePage(page, pageSize)
	list, total, err := s.repo.ListResources(resType, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(list))
	for _, r := range list {
		ids = append(ids, r.AuthorID)
	}
	authors := s.authors(ids)
	for i := range list {
		list[i].Author = attachAuthor(authors, list[i].AuthorID)
	}
	return &model.PageResult[model.KnowledgeResource]{List: nonNil(list), Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *CommunityService) GetResource(id string) (*model.KnowledgeResource, error) {
	r, err := s.repo.FindResource(id)
	if err != nil || r.Status != "published" {
		return nil, apperr.NotFound("资源不存在")
	}
	r.Views++
	_ = s.repo.SaveResource(r)
	authors := s.authors([]string{r.AuthorID})
	r.Author = attachAuthor(authors, r.AuthorID)
	return r, nil
}

func (s *CommunityService) CreateResource(userID string, req model.ResourceRequest) (*model.KnowledgeResource, error) {
	req.Title = strings.TrimSpace(req.Title)
	valid := false
	for _, t := range model.ResourceTypes {
		if t == req.Type {
			valid = true
		}
	}
	if !valid {
		return nil, apperr.BadRequest("资源类型无效")
	}
	if req.Title == "" || (strings.TrimSpace(req.Content) == "" && strings.TrimSpace(req.URL) == "") {
		return nil, apperr.BadRequest("请填写标题，并提供正文或链接")
	}
	if req.Summary == "" {
		req.Summary = excerpt(req.Content, 120)
	}
	r := &model.KnowledgeResource{
		ID: genID("res"), Type: req.Type, Title: req.Title, Summary: req.Summary, Content: req.Content, URL: req.URL,
		AuthorID: userID, SourceType: "creator", KnowledgeIDs: toJSON(nonNil(req.KnowledgeIDs)),
		SkillIDs: toJSON(nonNil(req.SkillIDs)), Tags: toJSON(nonNil(req.Tags)), Status: "published", CreatedAt: time.Now().UnixMilli(),
	}
	return r, s.repo.SaveResource(r)
}

func (s *CommunityService) DeleteResource(id string) error { return s.repo.DeleteResource(id) }

// ResourcesForKnowledge 与知识点关联的知识资产。
func (s *CommunityService) ResourcesForKnowledge(knowledgeID string, limit int) []model.KnowledgeResource {
	list, _, _ := s.repo.ListResources("", "", 1, 200)
	out := []model.KnowledgeResource{}
	for _, r := range list {
		for _, kid := range parseStrings(r.KnowledgeIDs) {
			if kid == knowledgeID {
				out = append(out, r)
				break
			}
		}
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func (s *CommunityService) Contributions(userID string) model.ContributionStats {
	answers, likes := s.repo.CountAnswersBy(userID)
	return model.ContributionStats{
		Posts: s.repo.CountPostsBy(userID), Answers: answers, AnswerLikes: likes, Resources: s.repo.CountResourcesBy(userID),
	}
}