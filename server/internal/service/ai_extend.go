package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/llm"
)

const maxHistoryTurns = 8

func buildLLMMessages(system string, history []model.ChatMessage, userPrompt string) []llm.Message {
	messages := []llm.Message{{Role: "system", Content: system}}
	for _, item := range trimHistory(history) {
		role := strings.TrimSpace(item.Role)
		content := strings.TrimSpace(item.Content)
		if role != "user" && role != "assistant" {
			continue
		}
		if content == "" {
			continue
		}
		messages = append(messages, llm.Message{Role: role, Content: content})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userPrompt})
	return messages
}

func trimHistory(history []model.ChatMessage) []model.ChatMessage {
	if len(history) == 0 {
		return nil
	}
	start := 0
	if len(history) > maxHistoryTurns*2 {
		start = len(history) - maxHistoryTurns*2
	}
	return history[start:]
}

// CareerInterview AI 职业访谈（多轮对话）
func (s *AIService) CareerInterview(ctx context.Context, message string, history []model.ChatMessage, onToken func(string)) (string, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return "", fmt.Errorf("请输入内容")
	}
	systemPrompt := `你是 IT 学习平台的职业选择顾问。通过 3-6 轮对话了解学习者的兴趣、基础、可投入时间与职业目标，语气亲切专业。
每次只问 1-2 个问题，不要一次问太多。当信息足够时，总结并推荐 1-3 个 IT 职业方向（如前端、后端、数据分析、大模型应用等），说明匹配理由。
使用中文 Markdown 输出。`

	if !s.settings.LLMClient().Enabled() {
		return "", fmt.Errorf("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}

	messages := []llm.Message{{Role: "system", Content: systemPrompt}}
	for _, item := range trimHistory(history) {
		if item.Role == "user" || item.Role == "assistant" {
			messages = append(messages, llm.Message{Role: item.Role, Content: item.Content})
		}
	}
	messages = append(messages, llm.Message{Role: "user", Content: message})
	return s.settings.LLMClient().StreamMessages(ctx, messages, onToken)
}

// CareerRecommend 根据背景推荐职业方向
func (s *AIService) CareerRecommend(ctx context.Context, req model.CareerRecommendRequest) (*model.CareerRecommendResult, error) {
	if !s.settings.LLMClient().Enabled() {
		return nil, fmt.Errorf("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}

	prompt := fmt.Sprintf(`根据以下学习者信息，推荐 2-3 个 IT 职业方向。
背景：%s
兴趣：%s
经验：%s
每周可投入：%d 小时

请严格输出 JSON（不要 markdown 代码块）：
{"summary":"一句话总结","recommendations":[{"careerId":"frontend|backend|data|llm","name":"职业名称","matchScore":85,"reason":"推荐理由"}]}`,
		req.Background, req.Interest, req.Experience, req.WeeklyHours)
	raw, err := s.settings.LLMClient().Complete(ctx, []llm.Message{
		{Role: "system", Content: "你是 IT 职业规划师，只输出合法 JSON。"},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, err
	}
	var result model.CareerRecommendResult
	if err := parseJSONResponse(raw, &result); err == nil && len(result.Recommendations) > 0 {
		return &result, nil
	}
	return nil, fmt.Errorf("大模型返回结果解析失败，请重试")
}

// GoalDecompose AI 学习目标分解
func (s *AIService) GoalDecompose(ctx context.Context, req model.GoalDecomposeRequest) (*model.GoalDecomposeResult, error) {
	if req.CareerName == "" {
		req.CareerName = "初级前端工程师"
	}
	if req.BaseLevel == "" {
		req.BaseLevel = "零基础"
	}
	if req.WeeklyHours <= 0 {
		req.WeeklyHours = 12
	}
	if req.DurationWeeks <= 0 {
		req.DurationWeeks = 16
	}

	if !s.settings.LLMClient().Enabled() {
		return nil, fmt.Errorf("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}

	prompt := fmt.Sprintf(`为学习者分解学习目标，输出 JSON（不要 markdown 代码块）：
职业：%s（%s）
当前基础：%s
每周投入：%d 小时
目标周期：%d 周

格式：
{"goalName":"目标名称","difficulty":"简单|中等|偏高","stages":[{"name":"阶段名","durationWeeks":4,"objectives":["目标"],"topics":["知识点"]}],"milestones":["里程碑"],"aiSummary":"总结","suggestions":["建议"]}`,
		req.CareerName, req.CareerID, req.BaseLevel, req.WeeklyHours, req.DurationWeeks)
	raw, err := s.settings.LLMClient().Complete(ctx, []llm.Message{
		{Role: "system", Content: "你是 IT 学习路径规划师，只输出合法 JSON，阶段数 3-5 个。"},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, err
	}
	var goalResult model.GoalDecomposeResult
	if err := parseJSONResponse(raw, &goalResult); err == nil && len(goalResult.Stages) > 0 {
		return &goalResult, nil
	}
	return nil, fmt.Errorf("大模型返回结果解析失败，请重试")
}

// LearningSuggest 生成个性化学习建议
func (s *AIService) LearningSuggest(ctx context.Context, req model.LearningSuggestRequest) (*model.LearningSuggestResult, error) {
	if !s.settings.LLMClient().Enabled() {
		return nil, fmt.Errorf("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}

	data, _ := json.Marshal(req)
	raw, err := s.settings.LLMClient().Complete(ctx, []llm.Message{
		{Role: "system", Content: "你是学习教练，根据学习数据给出 2-4 条简短中文建议。输出 JSON：{\"suggestions\":[\"建议1\",\"建议2\"]}"},
		{Role: "user", Content: string(data)},
	})
	if err != nil {
		return nil, err
	}
	var suggestResult model.LearningSuggestResult
	if err := parseJSONResponse(raw, &suggestResult); err == nil && len(suggestResult.Suggestions) > 0 {
		return &suggestResult, nil
	}
	return nil, fmt.Errorf("大模型返回结果解析失败，请重试")
}

func parseJSONResponse(raw string, dest interface{}) error {
	raw = strings.TrimSpace(raw)
	raw = regexp.MustCompile("(?s)^```(?:json)?\\s*").ReplaceAllString(raw, "")
	raw = regexp.MustCompile("(?s)\\s*```$").ReplaceAllString(raw, "")
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		raw = raw[start : end+1]
	}
	return json.Unmarshal([]byte(raw), dest)
}
