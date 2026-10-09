package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/llm"
	"github.com/cng1985/ai-learning-server/pkg/pipeline"
)

// learningState 在学习 Pipeline 各 Stage 之间传递的上下文。
type learningState struct {
	UserID   string
	Question string
	OnToken  func(string)

	User    *model.User
	Goal    *model.UserGoal
	Studied int
	Out     model.KernelOutput
}

// KernelService AI Learning Kernel：以 Pipeline 编排"用户上下文 → 职业分析 → 技能差距 → 知识推荐 → 学习规划 → AI 回答 → 评估 → 画像更新"。
type KernelService struct {
	growth   *GrowthService
	settings *SettingsService
	users    *repository.UserRepo
	repo     *repository.KernelRepo
	learning *pipeline.Pipeline[*learningState]
}

func NewKernelService(growth *GrowthService, settings *SettingsService, users *repository.UserRepo, repo *repository.KernelRepo) *KernelService {
	s := &KernelService{growth: growth, settings: settings, users: users, repo: repo}
	s.learning = pipeline.New[*learningState]("learning",
		stage("UserContextStage", "用户上下文", s.userContext),
		stage("CareerAnalysisStage", "职业分析", s.careerAnalysis),
		stage("SkillGapStage", "技能差距", s.skillGap),
		stage("KnowledgeRecommendStage", "知识推荐", s.knowledgeRecommend),
		stage("LearningPlanStage", "学习规划", s.learningPlan),
		stage("AIResponseStage", "AI 回答", s.aiResponse),
		stage("EvaluationStage", "规划评估", s.evaluation),
		stage("ProfileUpdateStage", "画像更新", s.profileUpdate),
	)
	return s
}

func stage(name, title string, fn func(context.Context, *learningState) (string, error)) pipeline.Stage[*learningState] {
	return pipeline.Func[*learningState]{StageName: name, StageTitle: title, Fn: fn}
}

// StageInfo Pipeline 结构描述。
type StageInfo struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (s *KernelService) Stages() []StageInfo {
	out := make([]StageInfo, 0, len(s.learning.Stages))
	for _, st := range s.learning.Stages {
		out = append(out, StageInfo{Name: st.Name(), Title: st.Title()})
	}
	return out
}

// RunLearning 执行学习 Pipeline 并持久化运行记录。
func (s *KernelService) RunLearning(ctx context.Context, userID, question string, onTrace func(pipeline.Trace), onToken func(string)) (*model.PipelineRun, error) {
	st := &learningState{UserID: userID, Question: strings.TrimSpace(question), OnToken: onToken}
	start := time.Now()
	traces, err := s.learning.Run(ctx, st, onTrace)
	run := &model.PipelineRun{
		ID: genID("run"), UserID: userID, Pipeline: s.learning.Name, Status: pipeline.StatusOK,
		Trace: toJSON(traces), Output: toJSON(st.Out), DurationMs: time.Since(start).Milliseconds(), CreatedAt: start.UnixMilli(),
	}
	if err != nil {
		run.Status = pipeline.StatusFailed
	}
	_ = s.repo.Save(run)
	return run, err
}

func (s *KernelService) Runs(userID string, limit int) ([]model.PipelineRun, error) {
	return s.repo.List(userID, limit)
}

func (s *KernelService) userContext(_ context.Context, st *learningState) (string, error) {
	u, err := s.users.FindByID(st.UserID)
	if err != nil {
		return "", fmt.Errorf("用户不存在")
	}
	st.User = u
	st.Goal = s.growth.Goal(st.UserID)
	states, _ := s.growth.KnowledgeStates(st.UserID)
	st.Studied = len(states)
	return fmt.Sprintf("%s · 已有 %d 个知识点学习记录", u.Nickname, st.Studied), nil
}

func (s *KernelService) careerAnalysis(_ context.Context, st *learningState) (string, error) {
	if st.Goal == nil {
		return "", pipeline.Skip("尚未设定职业目标，将按全部技能给出通用建议")
	}
	_, role, err := s.growth.catalog.RoleFull(st.Goal.RoleID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("目标岗位「%s」，包含 %d 项能力，每周投入 %d 小时、目标 %d 周", role.Name, len(role.Capabilities), st.Goal.WeeklyHours, st.Goal.TargetWeeks), nil
}

func (s *KernelService) skillGap(_ context.Context, st *learningState) (string, error) {
	// 先按最新知识状态（含遗忘衰减）重算技能，保证差距分析基于当前真实水平
	if err := s.growth.RecomputeSkills(st.UserID, nil); err != nil {
		return "", err
	}
	if st.Goal == nil {
		return "", pipeline.Skip("无目标岗位，跳过差距分析")
	}
	gap, err := s.growth.Gap(st.UserID, st.Goal.RoleID)
	if err != nil {
		return "", err
	}
	st.Out.Gap = gap
	return fmt.Sprintf("岗位达成度 %d%%，%d/%d 项技能达标，%d 项待提升", gap.Readiness, gap.MetCount, gap.TotalCount, len(gap.Gaps)), nil
}

func (s *KernelService) knowledgeRecommend(_ context.Context, st *learningState) (string, error) {
	recs, err := s.growth.Recommend(st.UserID, st.Out.Gap, 12)
	if err != nil {
		return "", err
	}
	st.Out.Recommendations = recs
	ready := 0
	for _, r := range recs {
		if r.Ready {
			ready++
		}
	}
	return fmt.Sprintf("按知识图谱前置关系推荐 %d 个知识点，其中 %d 个可立即开始", len(recs), ready), nil
}

func (s *KernelService) learningPlan(_ context.Context, st *learningState) (string, error) {
	st.Out.Plan = s.growth.BuildPlan(st.UserID, st.Goal, st.Out.Gap, st.Out.Recommendations)
	return st.Out.Plan.Summary, nil
}

func (s *KernelService) aiResponse(ctx context.Context, st *learningState) (string, error) {
	brief := s.planBrief(st)
	client := s.settings.LLMClient()
	if client == nil || !client.Enabled() {
		st.Out.Response, st.Out.ResponseSource = s.templateResponse(st), "rule"
		if st.OnToken != nil {
			st.OnToken(st.Out.Response)
		}
		return "", pipeline.Skip("大模型未配置，已使用规则模板生成学习建议")
	}
	question := st.Question
	if question == "" {
		question = "请根据我的情况给出本周学习建议。"
	}
	messages := []llm.Message{
		{Role: "system", Content: "你是 AI 学习内核中的 Learning Agent。基于系统给出的结构化分析（技能差距、知识推荐、学习计划）回答学习者，语气鼓励且具体，使用中文 Markdown，控制在 300 字以内。不要编造数据。"},
		{Role: "user", Content: fmt.Sprintf("学习者：%s\n%s\n\n学习者的问题：%s", st.User.Nickname, brief, question)},
	}
	text, err := client.StreamMessages(ctx, messages, st.OnToken)
	if err != nil {
		st.Out.Response, st.Out.ResponseSource = s.templateResponse(st), "rule"
		return "", pipeline.Skip("大模型调用失败，已回退到规则模板：" + err.Error())
	}
	st.Out.Response, st.Out.ResponseSource = text, "ai"
	return fmt.Sprintf("Learning Agent 已生成个性化建议（%d 字）", len([]rune(text))), nil
}

func (s *KernelService) planBrief(st *learningState) string {
	data := map[string]interface{}{}
	if st.Out.Gap != nil {
		gaps := []string{}
		for i, g := range st.Out.Gap.Gaps {
			if i >= 6 {
				break
			}
			gaps = append(gaps, fmt.Sprintf("%s L%d→L%d", g.SkillName, g.Current, g.Required))
		}
		data["targetRole"] = st.Out.Gap.Role.Name
		data["readiness"] = st.Out.Gap.Readiness
		data["skillGaps"] = gaps
	}
	recs := []string{}
	for i, r := range st.Out.Recommendations {
		if i >= 6 {
			break
		}
		recs = append(recs, r.Knowledge.Name)
	}
	data["nextKnowledge"] = recs
	if st.Out.Plan != nil {
		data["plan"] = st.Out.Plan.Summary
	}
	b, _ := json.Marshal(data)
	return "系统分析：" + string(b)
}

func (s *KernelService) templateResponse(st *learningState) string {
	var b strings.Builder
	if st.Out.Gap != nil {
		fmt.Fprintf(&b, "你距离「%s」的岗位要求已完成 **%d%%**。", st.Out.Gap.Role.Name, st.Out.Gap.Readiness)
		if len(st.Out.Gap.Gaps) > 0 {
			g := st.Out.Gap.Gaps[0]
			fmt.Fprintf(&b, "当前最关键的差距是 **%s**（L%d → L%d）。", g.SkillName, g.Current, g.Required)
		}
	} else {
		b.WriteString("你还没有设定职业目标，建议先在「职业与能力」中选择目标岗位，系统会据此分析差距。")
	}
	if len(st.Out.Recommendations) > 0 {
		b.WriteString("\n\n**本周建议优先学习：**\n")
		for i, r := range st.Out.Recommendations {
			if i >= 3 {
				break
			}
			fmt.Fprintf(&b, "%d. %s —— %s\n", i+1, r.Knowledge.Name, r.Reason)
		}
	}
	b.WriteString("\n学完后记得做练习：只有练习正确率与项目证据才能把\"知道\"变成\"会做\"。")
	return b.String()
}

func (s *KernelService) evaluation(_ context.Context, st *learningState) (string, error) {
	ev := model.KernelEvaluation{Coverage: 100, Feasibility: 100, Notes: []string{}}
	if st.Out.Gap != nil && len(st.Out.Gap.Gaps) > 0 {
		covered := map[string]bool{}
		for _, r := range st.Out.Recommendations {
			covered[r.SkillID] = true
		}
		hit := 0
		for _, g := range st.Out.Gap.Gaps {
			if covered[g.SkillID] {
				hit++
			}
		}
		ev.Coverage = hit * 100 / len(st.Out.Gap.Gaps)
		if ev.Coverage < 100 {
			ev.Notes = append(ev.Notes, "部分差距技能需要通过项目实践或新增知识点覆盖")
		}
	}
	if st.Goal != nil && st.Out.Plan != nil && len(st.Out.Plan.Weeks) > 0 {
		weeks := len(st.Out.Plan.Weeks)
		if weeks > st.Goal.TargetWeeks {
			ev.Feasibility = st.Goal.TargetWeeks * 100 / weeks
			ev.Notes = append(ev.Notes, fmt.Sprintf("计划需要 %d 周，超过目标周期 %d 周，建议增加每周投入", weeks, st.Goal.TargetWeeks))
		}
	}
	if len(st.Out.Recommendations) == 0 {
		ev.Notes = append(ev.Notes, "当前知识点均已掌握，下一步重点是项目实践与真实任务")
	}
	st.Out.Evaluation = ev
	return fmt.Sprintf("差距覆盖率 %d%%，计划可行性 %d%%", ev.Coverage, ev.Feasibility), nil
}

func (s *KernelService) profileUpdate(_ context.Context, st *learningState) (string, error) {
	states, _ := s.growth.SkillStates(st.UserID)
	updated := 0
	for _, v := range states {
		if v.Level > 0 {
			updated++
		}
	}
	st.Out.Profile = model.ProfileDelta{SkillsUpdated: updated}
	if st.Out.Gap != nil {
		st.Out.Profile.Readiness = st.Out.Gap.Readiness
	}
	return fmt.Sprintf("人才画像已同步：%d 项技能有等级记录", updated), nil
}
