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

func (s *AIService) Chat(ctx context.Context, req model.ChatRequest, onToken func(string)) (*model.ChatResult, error) {
	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, fmt.Errorf("问题不能为空")
	}

	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode == "" {
		mode = "rag"
	}

	var sources []model.AISource
	var systemPrompt string
	var userPrompt string

	if mode == "chat" {
		systemPrompt = "你是 ChatGPT 风格的智能助手，回答准确、简洁，使用中文 Markdown。可以进行通用对话、代码解释、文案撰写等任务。"
		userPrompt = question
	} else {
		_ = s.knowledge.EnsureIndexed(ctx)
		results, err := s.knowledge.Search(ctx, question, 0)
		if err != nil {
			return nil, err
		}
		sources = ToAISources(results)
		contextText := BuildContextFromResults(results)
		systemPrompt = "你是 AI 学习助手，基于提供的课程知识库回答用户问题。回答要准确、简洁，使用中文 Markdown。如果知识库中没有相关信息，请诚实说明，并给出学习建议。"
		userPrompt = fmt.Sprintf("参考知识库：\n%s\n\n用户问题：%s", contextText, question)
	}

	client, resolved, err := s.settings.LLMClientFor(req.VirtualModel)
	if err == nil && client != nil && client.Enabled() {
		messages := buildLLMMessages(systemPrompt, req.History, userPrompt)
		full, err := client.StreamMessages(ctx, messages, onToken)
		if err != nil {
			return nil, err
		}
		result := &model.ChatResult{Text: full, Sources: sources}
		if resolved != nil {
			result.Model = resolved.ModelCode
			result.Provider = resolved.ProviderCode
			result.VirtualModel = resolved.VirtualModelCode
			result.CanonicalModel = resolved.CanonicalModelCode
		}
		return result, nil
	}

	if mode == "chat" {
		answer := localChatAnswer(question)
		for _, ch := range answer {
			if onToken != nil {
				onToken(string(ch))
			}
			time.Sleep(15 * time.Millisecond)
		}
		return &model.ChatResult{Text: answer, Sources: sources}, nil
	}

	results, _ := s.knowledge.Search(ctx, question, 0)
	answer := localAnswerFromResults(question, results)
	for _, ch := range answer {
		if onToken != nil {
			onToken(string(ch))
		}
		time.Sleep(15 * time.Millisecond)
	}
	return &model.ChatResult{Text: answer, Sources: sources}, nil
}

func localChatAnswer(question string) string {
	return fmt.Sprintf(
		"当前处于**本地模式**（未配置 API Key 或模型路由不可用）。\n\n"+
			"你问的是：「%s」\n\n"+
			"请在管理端 **AI 大模型配置** 中为厂商填写 API Key，或在 **系统设置** 中配置 `LLM_API_KEY` 后重试。",
		question,
	)
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
		"enabled":       cfg.Enabled,
		"model":         cfg.Model,
		"baseUrl":       cfg.BaseURL,
		"source":        view.LLM.Source,
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
