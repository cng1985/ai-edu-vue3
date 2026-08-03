package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
)

type AIService struct {
	courses   *repository.CourseRepo
	settings  *SettingsService
	knowledge *KnowledgeService
}

func NewAIService(courses *repository.CourseRepo, settings *SettingsService, knowledge *KnowledgeService) *AIService {
	return &AIService{
		courses:   courses,
		settings:  settings,
		knowledge: knowledge,
	}
}

func (s *AIService) Chat(ctx context.Context, question string, history []model.ChatMessage, onToken func(string)) (*model.ChatResult, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("问题不能为空")
	}

	_ = s.knowledge.EnsureIndexed(ctx)

	results, err := s.knowledge.Search(ctx, question, 0)
	if err != nil {
		return nil, err
	}
	sources := ToAISources(results)
	contextText := BuildContextFromResults(results)

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

	answer := localAnswerFromResults(question, results)
	for _, ch := range answer {
		if onToken != nil {
			onToken(string(ch))
		}
		time.Sleep(15 * time.Millisecond)
	}
	return &model.ChatResult{Text: answer, Sources: sources}, nil
}

func localAnswerFromResults(question string, results []model.KnowledgeSearchResult) string {
	if len(results) == 0 {
		return "抱歉，我在知识库中没有找到与「" + question + "」直接相关的内容。你可以尝试换个问法，或前往课程页面系统学习。"
	}
	var b strings.Builder
	b.WriteString("根据课程内容，关于「")
	b.WriteString(question)
	b.WriteString("」的解答如下：\n\n")
	for i, r := range results {
		excerpt := r.Chunk.Text
		if len(excerpt) > 200 {
			excerpt = excerpt[:200] + "…"
		}
		b.WriteString(fmt.Sprintf("%d. **%s**（%s）\n%s\n\n", i+1, r.Chunk.Heading, r.Chunk.ChapterTitle, excerpt))
	}
	b.WriteString("---\n*以上回答基于向量知识库检索生成。配置 LLM_API_KEY 后可接入真实大模型。*")
	return b.String()
}

// ConfigInfo 返回 AI 配置状态（不含密钥）
func (s *AIService) ConfigInfo() map[string]interface{} {
	cfg := s.settings.LLMConfig()
	view := s.settings.GetView()
	kbStatus, _ := s.knowledge.Status()
	out := map[string]interface{}{
		"enabled": cfg.Enabled,
		"model":   cfg.Model,
		"baseUrl": cfg.BaseURL,
		"source":  view.LLM.Source,
		"knowledgeBase": kbStatus,
	}
	if s.settings.DefaultVirtualModel() != "" {
		out["defaultVirtualModel"] = s.settings.DefaultVirtualModel()
		resolved, err := s.settings.ResolvedLLM()
		if err == nil && resolved != nil {
			out["canonicalModel"] = resolved.CanonicalModelCode
			out["provider"] = resolved.ProviderCode
		}
	}
	return out
}
