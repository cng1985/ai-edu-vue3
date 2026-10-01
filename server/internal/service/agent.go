package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
)

type agentDef struct {
	info    model.AgentInfo
	persona string
	context func(ctx context.Context, userID, role, message string) (string, []model.AISource)
}

// AgentService AI Agent 体系：每个 Agent 拥有独立职责、人设与上下文工具。
type AgentService struct {
	growth    *GrowthService
	projects  *ProjectService
	market    *MarketService
	knowledge *KnowledgeService
	settings  *SettingsService
	catalog   *CatalogService
	agents    []agentDef
}

func NewAgentService(growth *GrowthService, projects *ProjectService, market *MarketService, knowledge *KnowledgeService,
	settings *SettingsService, catalog *CatalogService) *AgentService {
	s := &AgentService{growth: growth, projects: projects, market: market, knowledge: knowledge, settings: settings, catalog: catalog}
	s.agents = []agentDef{
		{
			info: model.AgentInfo{
				Code: "career", Name: "Career Agent", Title: "职业规划师", Icon: "compass", Color: "#6b5cff",
				Description:      "理解你的背景与兴趣，分析岗位能力模型，规划职业成长路线。",
				Responsibilities: []string{"职业规划", "能力分析", "成长路线"},
				Starters:         []string{"我适合哪个职业方向？", "我离目标岗位还差什么？", "帮我规划未来半年的成长路线"},
			},
			persona: "你是 Career Agent（职业规划师），负责职业规划、能力分析与成长路线。基于平台的职业-岗位-能力-技能模型给出建议，每次只问 1-2 个关键问题。",
			context: s.careerContext,
		},
		{
			info: model.AgentInfo{
				Code: "learning", Name: "Learning Agent", Title: "学习教练", Icon: "zap", Color: "#0fb981",
				Description:      "制定学习计划、讲解知识、监督学习节奏与间隔复习。",
				Responsibilities: []string{"学习计划", "知识解释", "学习监督"},
				Starters:         []string{"今天我该学什么？", "帮我安排本周学习计划", "我最近的学习状态怎么样？"},
			},
			persona: "你是 Learning Agent（学习教练），负责学习计划、知识解释与学习监督。用掌握度与复习计划数据说话，建议要具体到知识点与时长。",
			context: s.learningContext,
		},
		{
			info: model.AgentInfo{
				Code: "knowledge", Name: "Knowledge Agent", Title: "知识专家", Icon: "brain", Color: "#2a8cf4",
				Description:      "基于知识图谱与课程知识库解答问题，梳理知识关联。",
				Responsibilities: []string{"知识理解", "知识关联", "知识生成"},
				Starters:         []string{"volatile 和 synchronized 有什么区别？", "RAG 的混合检索是怎么工作的？", "AQS 的前置知识有哪些？"},
			},
			persona: "你是 Knowledge Agent（知识专家），负责知识理解、知识关联与知识生成。回答时说明前置知识与关联知识，引用知识库内容时保持准确，不知道就说明。",
			context: s.knowledgeContext,
		},
		{
			info: model.AgentInfo{
				Code: "project", Name: "Project Agent", Title: "项目导师", Icon: "layers", Color: "#f59e0b",
				Description:      "指导项目实践，拆解任务，给出技术方案建议。",
				Responsibilities: []string{"项目指导", "任务拆解", "技术方案"},
				Starters:         []string{"推荐一个适合我练手的项目任务", "订单系统的库存扣减怎么设计？", "帮我拆解 RAG 知识库项目"},
			},
			persona: "你是 Project Agent（项目导师），负责项目指导、任务拆解与技术方案。给出可落地的步骤、关键设计点与验收标准，引导学习者自己动手而不是直接给完整答案。",
			context: s.projectContext,
		},
		{
			info: model.AgentInfo{
				Code: "review", Name: "Review Agent", Title: "能力评审官", Icon: "clipboard", Color: "#ef4d63",
				Description:      "分析代码与项目成果，评价项目质量，评估能力等级。",
				Responsibilities: []string{"代码分析", "项目评价", "能力评估"},
				Starters:         []string{"帮我 review 这段代码", "我最近的项目提交有什么问题？", "我的技能等级为什么还是 L2？"},
			},
			persona: "你是 Review Agent（能力评审官），负责代码分析、项目评价与能力评估。评审要客观、指出具体问题与改进方向；解释技能等级时依据'知识状态 + 项目证据 + 任务表现'模型。",
			context: s.reviewContext,
		},
		{
			info: model.AgentInfo{
				Code: "opportunity", Name: "Opportunity Agent", Title: "机会发现官", Icon: "target", Color: "#14b8a6",
				Description:      "匹配真实 IT 任务，发现成长与收益机会。",
				Responsibilities: []string{"任务匹配", "人才推荐", "机会发现"},
				Starters:         []string{"有哪些任务适合我？", "我还差什么才能接这个任务？", "怎样提高我的任务匹配度？"},
			},
			persona: "你是 Opportunity Agent（机会发现官），负责任务匹配、人才推荐与机会发现。基于技能要求与人才画像的匹配数据给建议，帮助学习者通过真实任务获得收益并持续成长。",
			context: s.opportunityContext,
		},
	}
	return s
}

func (s *AgentService) List() []model.AgentInfo {
	out := make([]model.AgentInfo, 0, len(s.agents))
	for _, a := range s.agents {
		out = append(out, a.info)
	}
	return out
}

func (s *AgentService) find(code string) (*agentDef, bool) {
	for i := range s.agents {
		if s.agents[i].info.Code == code {
			return &s.agents[i], true
		}
	}
	return nil, false
}

// AgentReply Agent 回复结果。
type AgentReply struct {
	Text    string           `json:"text"`
	Sources []model.AISource `json:"sources"`
	Context string           `json:"context"`
}

// Chat 与指定 Agent 对话：注入该 Agent 职责相关的个人成长上下文后流式生成。
func (s *AgentService) Chat(ctx context.Context, userID, role, code string, req model.AgentChatRequest, onToken func(string)) (*AgentReply, error) {
	agent, ok := s.find(code)
	if !ok {
		return nil, apperr.NotFound("Agent 不存在")
	}
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		return nil, apperr.BadRequest("请输入内容")
	}
	client := s.settings.LLMClient()
	if client == nil || !client.Enabled() {
		return nil, apperr.BadRequest("大模型未配置，请在管理端「大模型配置」中为厂商填写 API Key")
	}
	contextText, sources := agent.context(ctx, userID, role, msg)
	system := agent.persona + "\n使用中文 Markdown 回答。以下是平台提供的学习者实时数据，仅在相关时引用：\n<context>\n" + contextText + "\n</context>"
	text, err := client.StreamMessages(ctx, buildLLMMessages(system, req.History, msg), onToken)
	if err != nil {
		return nil, err
	}
	return &AgentReply{Text: text, Sources: nonNil(sources), Context: contextText}, nil
}

// ContextPreview 返回 Agent 将使用的上下文（便于学习者理解 AI 依据）。
func (s *AgentService) ContextPreview(ctx context.Context, userID, role, code, message string) (string, error) {
	agent, ok := s.find(code)
	if !ok {
		return "", apperr.NotFound("Agent 不存在")
	}
	text, _ := agent.context(ctx, userID, role, message)
	return text, nil
}

func (s *AgentService) goalLine(userID string) string {
	gap, err := s.growth.Gap(userID, "")
	if err != nil || gap == nil {
		return "职业目标：未设定"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "职业目标：%s / %s，岗位达成度 %d%%（%d/%d 项技能达标）\n", gap.Career.Name, gap.Role.Name, gap.Readiness, gap.MetCount, gap.TotalCount)
	for i, g := range gap.Gaps {
		if i >= 6 {
			break
		}
		fmt.Fprintf(&b, "- 差距：%s（%s）L%d → L%d\n", g.SkillName, g.CapabilityName, g.Current, g.Required)
	}
	return b.String()
}

func (s *AgentService) skillLines(userID string, limit int) string {
	states, err := s.growth.SkillStates(userID)
	if err != nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("技能状态（等级｜知识/项目/任务得分）：\n")
	n := 0
	for _, st := range states {
		if st.Level == 0 && st.KnowledgeScore == 0 {
			continue
		}
		fmt.Fprintf(&b, "- %s L%d %s｜%.0f%%/%.0f%%/%.0f%%\n", st.Skill.Name, st.Level, st.LevelName, st.KnowledgeScore*100, st.ProjectScore*100, st.TaskScore*100)
		n++
		if n >= limit {
			break
		}
	}
	if n == 0 {
		b.WriteString("- 暂无技能记录\n")
	}
	return b.String()
}

func (s *AgentService) careerContext(_ context.Context, userID, _, _ string) (string, []model.AISource) {
	var b strings.Builder
	b.WriteString(s.goalLine(userID))
	if careers, err := s.catalog.ListCareers(); err == nil {
		b.WriteString("平台职业体系：\n")
		for _, c := range careers {
			names := []string{}
			for _, r := range c.Roles {
				names = append(names, r.Name)
			}
			fmt.Fprintf(&b, "- %s（需求%s，薪资%s）：%s\n", c.Name, c.Demand, c.SalaryRange, strings.Join(names, "、"))
		}
	}
	b.WriteString(s.skillLines(userID, 8))
	return b.String(), nil
}

func (s *AgentService) learningContext(_ context.Context, userID, _, _ string) (string, []model.AISource) {
	var b strings.Builder
	b.WriteString(s.goalLine(userID))
	gap, _ := s.growth.Gap(userID, "")
	if recs, err := s.growth.Recommend(userID, gap, 8); err == nil {
		b.WriteString("推荐学习（按前置关系排序）：\n")
		for _, r := range recs {
			state := "可开始"
			if !r.Ready {
				state = "需先完成前置"
			}
			fmt.Fprintf(&b, "- %s（%s，掌握度 %.0f%%，%s，约 %d 分钟）\n", r.Knowledge.Name, r.SkillName, r.Mastery*100, state, r.Knowledge.EstimatedMinutes)
		}
	}
	due := s.growth.DueReviews(userID, 6)
	if len(due) > 0 {
		b.WriteString("到期复习：\n")
		for _, d := range due {
			fmt.Fprintf(&b, "- %s（当前掌握度 %.0f%%）\n", d.Knowledge.Name, d.EffectiveMastery*100)
		}
	}
	f := s.growth.Flywheel(userID, 0)
	fmt.Fprintf(&b, "累计学习 %d 次、%d 分钟，已掌握 %d 个知识点\n", f.LearningEvents, f.LearningMinutes, f.MasteredCount)
	return b.String(), nil
}

func (s *AgentService) knowledgeContext(ctx context.Context, userID, _, message string) (string, []model.AISource) {
	var b strings.Builder
	var sources []model.AISource
	if c, err := s.catalog.Load(); err == nil {
		mastery := s.growth.effectiveMasteryMap(userID)
		lower := strings.ToLower(message)
		hits := 0
		for _, k := range c.KnowledgeList {
			if !strings.Contains(lower, strings.ToLower(k.Name)) {
				continue
			}
			prereqs := []string{}
			for _, p := range c.Prereqs[k.ID] {
				prereqs = append(prereqs, c.Brief(p).Name)
			}
			fmt.Fprintf(&b, "知识点「%s」：%s；前置：%s；学习者掌握度 %.0f%%\n", k.Name, k.Summary, strings.Join(prereqs, "、"), mastery[k.ID]*100)
			hits++
			if hits >= 4 {
				break
			}
		}
	}
	_ = s.knowledge.EnsureIndexed(ctx)
	if results, err := s.knowledge.Search(ctx, message, 0); err == nil && len(results) > 0 {
		b.WriteString("课程知识库检索结果：\n")
		b.WriteString(BuildContextFromResults(results))
		sources = ToAISources(results)
	}
	if b.Len() == 0 {
		b.WriteString("知识库中没有直接相关的内容。")
	}
	return b.String(), sources
}

func (s *AgentService) projectContext(_ context.Context, userID, _, _ string) (string, []model.AISource) {
	var b strings.Builder
	b.WriteString(s.goalLine(userID))
	gap, _ := s.growth.Gap(userID, "")
	tasks := s.growth.PracticeTasks(userID, gap, 5)
	if len(tasks) > 0 {
		b.WriteString("推荐练习任务：\n")
		for _, t := range tasks {
			fmt.Fprintf(&b, "- %s · %s（匹配度 %d%%）\n", t.ProjectTitle, t.TaskTitle, t.Match.Score)
		}
	}
	if projects, err := s.projects.List(userID, "published"); err == nil {
		b.WriteString("平台项目：\n")
		for _, p := range projects {
			titles := []string{}
			for _, t := range p.Tasks {
				titles = append(titles, t.Title)
			}
			fmt.Fprintf(&b, "- %s（完成 %d%%）：%s\n", p.Title, p.Progress, strings.Join(titles, " → "))
		}
	}
	return b.String(), nil
}

func (s *AgentService) reviewContext(_ context.Context, userID, _, _ string) (string, []model.AISource) {
	var b strings.Builder
	b.WriteString(s.skillLines(userID, 10))
	b.WriteString("技能等级规则：仅有知识最高 L2；L3 需项目证据≥60 分；L4 需高质量项目(≥80)或真实任务(≥75)；L5 需多项证据且真实任务优秀。\n")
	subs := s.projects.Submissions(userID, 5)
	if len(subs) > 0 {
		b.WriteString("最近项目提交：\n")
		for _, sub := range subs {
			fmt.Fprintf(&b, "- %s · %s：%d 分（%s 评审）\n", sub.ProjectTitle, sub.TaskTitle, sub.Score, sub.Evaluator)
		}
	}
	return b.String(), nil
}

func (s *AgentService) opportunityContext(_ context.Context, userID, role, _ string) (string, []model.AISource) {
	var b strings.Builder
	b.WriteString(s.skillLines(userID, 8))
	list := s.market.TopMatches(userID, role, 6)
	if len(list) > 0 {
		b.WriteString("开放任务匹配：\n")
		for _, o := range list {
			reqs := []string{}
			for _, it := range o.Match.Items {
				reqs = append(reqs, fmt.Sprintf("%s L%d(当前L%d)", it.SkillName, it.Required, it.Current))
			}
			fmt.Fprintf(&b, "- %s · %s，预算 %.0f 元，匹配度 %d%%：%s\n", o.Company, o.Title, o.Budget, o.Match.Score, strings.Join(reqs, "、"))
		}
	} else {
		b.WriteString("当前没有可匹配的开放任务。\n")
	}
	return b.String(), nil
}
