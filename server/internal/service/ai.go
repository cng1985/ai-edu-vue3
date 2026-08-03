package service

import (
	"context"
	"fmt"
	"strings"

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
	if err != nil || client == nil || !client.Enabled() {
		return nil, fmt.Errorf("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}

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
