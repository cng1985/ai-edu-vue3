package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
)

type AIService struct {
	courses  *repository.CourseRepo
	settings *SettingsService
}

func NewAIService(courses *repository.CourseRepo, settings *SettingsService) *AIService {
	return &AIService{
		courses:  courses,
		settings: settings,
	}
}

type knowledgeChunk struct {
	courseID, courseTitle, chapterID, chapterTitle, heading, text string
}

func (s *AIService) Chat(ctx context.Context, question string, history []model.ChatMessage, onToken func(string)) (*model.ChatResult, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("问题不能为空")
	}
	chunks, err := s.buildKnowledgeBase()
	if err != nil {
		return nil, err
	}
	matched, sources := retrieve(chunks, question, 3)
	contextText := buildContext(matched)
	systemPrompt := "你是 AI 学习助手，基于提供的课程知识库回答用户问题。回答要准确、简洁，使用中文 Markdown。如果知识库中没有相关信息，请诚实说明，并给出学习建议。"
	userPrompt := fmt.Sprintf("参考知识库：\n%s\n\n用户问题：%s", contextText, question)

	if s.settings.LLMClient().Enabled() {
		messages := buildLLMMessages(systemPrompt, history, userPrompt)
		full, err := s.settings.LLMClient().StreamMessages(ctx, messages, onToken)
		if err != nil {
			return nil, err
		}
		return &model.ChatResult{Text: full, Sources: sources}, nil
	}
	// 无 LLM 配置时本地模拟
	answer := localAnswer(question, matched)
	for _, ch := range answer {
		if onToken != nil {
			onToken(string(ch))
		}
		time.Sleep(15 * time.Millisecond)
	}
	return &model.ChatResult{Text: answer, Sources: sources}, nil
}

func (s *AIService) buildKnowledgeBase() ([]knowledgeChunk, error) {
	courses, err := s.courses.ListPublished()
	if err != nil {
		return nil, err
	}
	var chunks []knowledgeChunk
	for _, course := range courses {
		for _, ch := range course.Chapters {
			for _, section := range splitByHeadings(ch.Content) {
				plain := stripMarkdown(section.body)
				if len(plain) < 30 {
					continue
				}
				chunks = append(chunks, knowledgeChunk{
					courseID: course.ID, courseTitle: course.Title,
					chapterID: ch.ID, chapterTitle: ch.Title,
					heading: section.heading, text: plain,
				})
			}
		}
	}
	return chunks, nil
}

type mdSection struct{ heading, body string }

func splitByHeadings(markdown string) []mdSection {
	lines := strings.Split(markdown, "\n")
	var sections []mdSection
	current := mdSection{}
	inCode := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCode = !inCode
		}
		if !inCode && regexp.MustCompile(`^#{1,3}\s`).MatchString(line) {
			if strings.TrimSpace(current.body) != "" {
				sections = append(sections, current)
			}
			current = mdSection{heading: regexp.MustCompile(`^#+\s*`).ReplaceAllString(line, ""), body: ""}
		} else {
			current.body += line + "\n"
		}
	}
	if strings.TrimSpace(current.body) != "" {
		sections = append(sections, current)
	}
	return sections
}

func stripMarkdown(text string) string {
	re := regexp.MustCompile("(?s)```.*?```")
	text = re.ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\|.*\|`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`[#>*` + "`" + `_[\]()!-]`).ReplaceAllString(text, " ")
	return strings.Join(strings.Fields(text), " ")
}

func retrieve(chunks []knowledgeChunk, query string, topK int) ([]knowledgeChunk, []model.AISource) {
	tokens := tokenize(query)
	type scored struct {
		chunk knowledgeChunk
		score float64
	}
	var results []scored
	for _, ch := range chunks {
		score := 0.0
		lower := strings.ToLower(ch.text + " " + ch.heading)
		for tok := range tokens {
			if strings.Contains(lower, tok) {
				score += 1
			}
		}
		if score > 0 {
			results = append(results, scored{chunk: ch, score: score})
		}
	}
	// 简单排序
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	if len(results) > topK {
		results = results[:topK]
	}
	var matched []knowledgeChunk
	var sources []model.AISource
	seen := map[string]bool{}
	for _, r := range results {
		matched = append(matched, r.chunk)
		key := r.chunk.courseID + ":" + r.chunk.chapterID
		if !seen[key] {
			seen[key] = true
			sources = append(sources, model.AISource{
				CourseID: r.chunk.courseID, CourseTitle: r.chunk.courseTitle,
				ChapterID: r.chunk.chapterID, ChapterTitle: r.chunk.chapterTitle,
			})
		}
	}
	return matched, sources
}

func tokenize(query string) map[string]bool {
	tokens := map[string]bool{}
	lower := strings.ToLower(query)
	for _, m := range regexp.MustCompile(`[a-z0-9-]{2,}`).FindAllString(lower, -1) {
		tokens[m] = true
	}
	runes := []rune(query)
	for i := 0; i < len(runes); i++ {
		for l := 2; l <= 4 && i+l <= len(runes); l++ {
			sub := string(runes[i : i+l])
			if isMostlyChinese(sub) {
				tokens[sub] = true
			}
		}
	}
	stop := map[string]bool{"什么": true, "怎么": true, "如何": true, "的": true, "了": true, "是": true}
	for k := range tokens {
		if stop[k] {
			delete(tokens, k)
		}
	}
	return tokens
}

func isMostlyChinese(s string) bool {
	n := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			n++
		}
	}
	return n >= len(s)/2
}

func buildContext(chunks []knowledgeChunk) string {
	var parts []string
	for _, ch := range chunks {
		parts = append(parts, fmt.Sprintf("【%s > %s > %s】\n%s", ch.courseTitle, ch.chapterTitle, ch.heading, ch.text))
	}
	if len(parts) == 0 {
		return "（未检索到相关知识）"
	}
	return strings.Join(parts, "\n\n")
}

func localAnswer(question string, chunks []knowledgeChunk) string {
	if len(chunks) == 0 {
		return "抱歉，我在知识库中没有找到与「" + question + "」直接相关的内容。你可以尝试换个问法，或前往课程页面系统学习。"
	}
	var b strings.Builder
	b.WriteString("根据课程内容，关于「")
	b.WriteString(question)
	b.WriteString("」的解答如下：\n\n")
	for i, ch := range chunks {
		excerpt := ch.text
		if len(excerpt) > 200 {
			excerpt = excerpt[:200] + "…"
		}
		b.WriteString(fmt.Sprintf("%d. **%s**（%s）\n%s\n\n", i+1, ch.heading, ch.chapterTitle, excerpt))
	}
	b.WriteString("---\n*以上回答基于课程知识库检索生成。配置 LLM_API_KEY 后可接入真实大模型。*")
	return b.String()
}

// ConfigInfo 返回 AI 配置状态（不含密钥）
func (s *AIService) ConfigInfo() map[string]interface{} {
	cfg := s.settings.LLMConfig()
	view := s.settings.GetView()
	return map[string]interface{}{
		"enabled": cfg.Enabled,
		"model":   cfg.Model,
		"baseUrl": cfg.BaseURL,
		"source":  view.LLM.Source,
	}
}

// suppress unused
var _ = json.Marshal
