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

	if s.llm.Enabled() {
		messages := []llm.Message{{Role: "system", Content: systemPrompt}}
		for _, item := range trimHistory(history) {
			if item.Role == "user" || item.Role == "assistant" {
				messages = append(messages, llm.Message{Role: item.Role, Content: item.Content})
			}
		}
		messages = append(messages, llm.Message{Role: "user", Content: message})
		return s.llm.StreamMessages(ctx, messages, onToken)
	}

	reply := localCareerInterview(message, len(history))
	for _, ch := range reply {
		if onToken != nil {
			onToken(string(ch))
		}
	}
	return reply, nil
}

func localCareerInterview(message string, historyLen int) string {
	lower := strings.ToLower(message)
	if historyLen == 0 {
		return "你好！我是你的职业规划助手。\n\n先聊聊你的情况：你目前是在校生、在职转行，还是准备求职？对哪类工作更感兴趣——做界面、写后端逻辑、分析数据，还是 AI 相关方向？"
	}
	if strings.Contains(lower, "前端") || strings.Contains(message, "界面") {
		return "听起来你对**前端开发**很有兴趣！\n\n再确认两点：\n1. 你目前 HTML/CSS/JS 基础如何？\n2. 每周大概能投入多少小时学习？\n\n如果基础较弱，建议目标定为「16 周成为初级前端工程师」，平台已有完整学习路径。"
	}
	return fmt.Sprintf("感谢分享！根据你提到的「%s」，我建议优先考虑 **Web 前端工程师**（匹配度约 88%%）。\n\n理由：入门路径清晰、作品可见、市场需求稳定。你可以在下一步确认目标条件，生成专属学习路径。\n\n*配置 LLM_API_KEY 后可获得更精准的多轮职业访谈。*", message)
}

// CareerRecommend 根据背景推荐职业方向
func (s *AIService) CareerRecommend(ctx context.Context, req model.CareerRecommendRequest) (*model.CareerRecommendResult, error) {
	if s.llm.Enabled() {
		prompt := fmt.Sprintf(`根据以下学习者信息，推荐 2-3 个 IT 职业方向。
背景：%s
兴趣：%s
经验：%s
每周可投入：%d 小时

请严格输出 JSON（不要 markdown 代码块）：
{"summary":"一句话总结","recommendations":[{"careerId":"frontend|backend|data|llm","name":"职业名称","matchScore":85,"reason":"推荐理由"}]}`,
			req.Background, req.Interest, req.Experience, req.WeeklyHours)
		raw, err := s.llm.Complete(ctx, []llm.Message{
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
	}
	return localCareerRecommend(req), nil
}

func localCareerRecommend(req model.CareerRecommendRequest) *model.CareerRecommendResult {
	interest := strings.ToLower(req.Interest + req.Background)
	recs := []model.CareerRecommendation{
		{CareerID: "frontend", Name: "Web 前端工程师", MatchScore: 88, Reason: "适合希望快速看到成果、通过作品展示能力的学习者"},
		{CareerID: "backend", Name: "Java/Go 后端工程师", MatchScore: 75, Reason: "适合逻辑性强、对系统架构感兴趣的学习者"},
		{CareerID: "llm", Name: "大模型应用工程师", MatchScore: 70, Reason: "适合对 AI、Prompt、RAG 感兴趣的学习者"},
	}
	if strings.Contains(interest, "数据") {
		recs[0], recs[2] = recs[2], recs[0]
		recs[0] = model.CareerRecommendation{CareerID: "data", Name: "数据分析师", MatchScore: 90, Reason: "与数据洞察兴趣高度匹配"}
	}
	return &model.CareerRecommendResult{
		Summary:         "综合你的兴趣与基础，以下方向值得优先考虑",
		Recommendations: recs,
	}
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

	if s.llm.Enabled() {
		prompt := fmt.Sprintf(`为学习者分解学习目标，输出 JSON（不要 markdown 代码块）：
职业：%s（%s）
当前基础：%s
每周投入：%d 小时
目标周期：%d 周

格式：
{"goalName":"目标名称","difficulty":"简单|中等|偏高","stages":[{"name":"阶段名","durationWeeks":4,"objectives":["目标"],"topics":["知识点"]}],"milestones":["里程碑"],"aiSummary":"总结","suggestions":["建议"]}`,
			req.CareerName, req.CareerID, req.BaseLevel, req.WeeklyHours, req.DurationWeeks)
		raw, err := s.llm.Complete(ctx, []llm.Message{
			{Role: "system", Content: "你是 IT 学习路径规划师，只输出合法 JSON，阶段数 3-5 个。"},
			{Role: "user", Content: prompt},
		})
		if err != nil {
			return nil, err
		}
		var result model.GoalDecomposeResult
		if err := parseJSONResponse(raw, &result); err == nil && len(result.Stages) > 0 {
			return &result, nil
		}
	}
	return localGoalDecompose(req), nil
}

func localGoalDecompose(req model.GoalDecomposeRequest) *model.GoalDecomposeResult {
	weeks := req.DurationWeeks
	return &model.GoalDecomposeResult{
		GoalName:   fmt.Sprintf("%d 周成为初级前端工程师", weeks),
		Difficulty: "中等",
		Stages: []model.LearningStage{
			{Name: "HTML/CSS 基础", DurationWeeks: weeks / 4, Objectives: []string{"能独立完成静态页面"}, Topics: []string{"语义化", "Flex 布局", "响应式"}},
			{Name: "JavaScript 核心", DurationWeeks: weeks / 4, Objectives: []string{"掌握 ES6+ 与异步"}, Topics: []string{"DOM", "Promise", "闭包"}},
			{Name: "前端框架", DurationWeeks: weeks / 4, Objectives: []string{"能用 Vue/React 开发页面"}, Topics: []string{"组件", "状态管理", "路由"}},
			{Name: "工程化与项目", DurationWeeks: weeks / 4, Objectives: []string{"完成可展示作品集"}, Topics: []string{"Vite", "Git", "项目实战"}},
		},
		Milestones: []string{"完成静态页面", "完成 Todo 项目", "完成后台管理系统", "通过模拟面试"},
		AISummary:  fmt.Sprintf("基于「%s」基础、每周 %d 小时，建议按 4 个阶段推进，重点保证项目产出。", req.BaseLevel, req.WeeklyHours),
		Suggestions: []string{
			"先完成平台微单元建立学习习惯",
			"每周至少完成 1 个小项目里程碑",
			"利用 AI 助手随时答疑",
		},
	}
}

// LearningSuggest 生成个性化学习建议
func (s *AIService) LearningSuggest(ctx context.Context, req model.LearningSuggestRequest) (*model.LearningSuggestResult, error) {
	if s.llm.Enabled() {
		data, _ := json.Marshal(req)
		raw, err := s.llm.Complete(ctx, []llm.Message{
			{Role: "system", Content: "你是学习教练，根据学习数据给出 2-4 条简短中文建议。输出 JSON：{\"suggestions\":[\"建议1\",\"建议2\"]}"},
			{Role: "user", Content: string(data)},
		})
		if err != nil {
			return nil, err
		}
		var result model.LearningSuggestResult
		if err := parseJSONResponse(raw, &result); err == nil && len(result.Suggestions) > 0 {
			return &result, nil
		}
	}
	return localLearningSuggest(req), nil
}

func localLearningSuggest(req model.LearningSuggestRequest) *model.LearningSuggestResult {
	var suggestions []string
	weakest := ""
	lowest := 101
	for _, c := range req.CompetencyProgress {
		if c.Progress < lowest {
			lowest = c.Progress
			weakest = c.Name
		}
	}
	if weakest != "" && lowest < 60 {
		suggestions = append(suggestions, fmt.Sprintf("你的「%s」掌握度仅 %d%%，建议先补强再推进下一阶段", weakest, lowest))
	}
	if req.NextMilestone != "" {
		suggestions = append(suggestions, fmt.Sprintf("下一里程碑「%s」即将到期，优先完成相关微单元", req.NextMilestone))
	}
	if req.Streak >= 3 {
		suggestions = append(suggestions, fmt.Sprintf("已连续学习 %d 天，保持节奏！", req.Streak))
	} else {
		suggestions = append(suggestions, "今天完成 2-3 个微单元，重新建立学习节奏")
	}
	suggestions = append(suggestions, "有疑问随时向 AI 学习助手提问")
	return &model.LearningSuggestResult{Suggestions: suggestions}
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
