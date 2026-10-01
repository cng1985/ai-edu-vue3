package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/growth"
)

// GrowthService 个人成长：职业目标、学习状态、技能评价、差距分析与知识推荐。
type GrowthService struct {
	catalog  *CatalogService
	repo     *repository.GrowthRepo
	projects *repository.ProjectRepo
	market   *repository.MarketRepo
	users    *repository.UserRepo
	kernel   *repository.KernelRepo
	now      func() int64
}

func NewGrowthService(catalog *CatalogService, repo *repository.GrowthRepo, projects *repository.ProjectRepo,
	market *repository.MarketRepo, users *repository.UserRepo, kernel *repository.KernelRepo) *GrowthService {
	return &GrowthService{
		catalog: catalog, repo: repo, projects: projects, market: market, users: users, kernel: kernel,
		now: func() int64 { return time.Now().UnixMilli() },
	}
}

// ---------- 职业目标 ----------

func (s *GrowthService) Goal(userID string) *model.UserGoal {
	g, err := s.repo.FindGoal(userID)
	if err != nil {
		return nil
	}
	return g
}

func (s *GrowthService) SetGoal(userID string, req model.GoalRequest) (*model.UserGoal, error) {
	career, role, err := s.catalog.RoleFull(req.RoleID)
	if err != nil {
		return nil, err
	}
	now := s.now()
	g := s.Goal(userID)
	if g == nil {
		g = &model.UserGoal{UserID: userID, CreatedAt: now}
	}
	g.CareerID, g.RoleID = career.ID, role.ID
	g.WeeklyHours = clampRange(req.WeeklyHours, 1, 80, 10)
	g.TargetWeeks = clampRange(req.TargetWeeks, 1, 104, 16)
	g.Motivation = strings.TrimSpace(req.Motivation)
	g.UpdatedAt = now
	if err := s.repo.SaveGoal(g); err != nil {
		return nil, err
	}
	return g, nil
}

func clampRange(v, lo, hi, def int) int {
	if v <= 0 {
		return def
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ---------- 学习事件与知识状态 ----------

func toEngineState(st model.StudentKnowledgeState) growth.KnowledgeState {
	return growth.KnowledgeState{
		StudyCount: st.StudyCount, PracticeCount: st.PracticeCount, CorrectCount: st.CorrectCount,
		Accuracy: st.Accuracy, Mastery: st.Mastery, Confidence: st.Confidence, Status: st.Status,
		LastStudiedAt: st.LastStudiedAt, NextReviewAt: st.NextReviewAt,
	}
}

func (s *GrowthService) stateView(c *Catalog, st model.StudentKnowledgeState, now int64) model.KnowledgeStateView {
	es := toEngineState(st)
	status := st.Status
	if status == "" {
		status = model.KnowledgeStatusNew
	}
	return model.KnowledgeStateView{
		Knowledge: c.Brief(st.KnowledgeID), StudyCount: st.StudyCount, PracticeCount: st.PracticeCount,
		Accuracy: st.Accuracy, Mastery: st.Mastery, EffectiveMastery: growth.EffectiveMastery(es, now),
		Confidence: st.Confidence, Status: status, Due: growth.IsDue(es, now),
		LastStudiedAt: st.LastStudiedAt, NextReviewAt: st.NextReviewAt,
	}
}

// RecordEvent 记录一次学习行为并实时更新知识状态与相关技能状态。
func (s *GrowthService) RecordEvent(userID string, req model.LearningEventRequest) (*model.KnowledgeStateView, error) {
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	view, err := s.applyEvent(c, userID, req)
	if err != nil {
		return nil, err
	}
	s.recomputeSkills(c, userID, c.KnowledgeSkills[req.KnowledgeID])
	return view, nil
}

func (s *GrowthService) applyEvent(c *Catalog, userID string, req model.LearningEventRequest) (*model.KnowledgeStateView, error) {
	k, ok := c.Knowledge[req.KnowledgeID]
	if !ok {
		return nil, apperr.NotFound("知识点不存在")
	}
	switch req.Type {
	case model.EventStudy, model.EventPractice, model.EventReview:
	default:
		return nil, apperr.BadRequest("学习事件类型无效")
	}
	if req.Type == model.EventPractice && req.Total <= 0 {
		return nil, apperr.BadRequest("练习事件需要题目数量")
	}
	now := s.now()
	st, err := s.repo.FindKnowledgeState(userID, k.ID)
	if err != nil {
		st = &model.StudentKnowledgeState{ID: genID("ks"), UserID: userID, KnowledgeID: k.ID}
	}
	next := growth.Apply(toEngineState(*st), growth.Event{Type: req.Type, Correct: req.Correct, Total: req.Total, At: now})
	st.StudyCount, st.PracticeCount, st.CorrectCount = next.StudyCount, next.PracticeCount, next.CorrectCount
	st.Accuracy, st.Mastery, st.Confidence, st.Status = next.Accuracy, next.Mastery, next.Confidence, next.Status
	st.LastStudiedAt, st.NextReviewAt, st.UpdatedAt = next.LastStudiedAt, next.NextReviewAt, now
	if err := s.repo.SaveKnowledgeState(st); err != nil {
		return nil, err
	}
	minutes := req.Minutes
	if minutes <= 0 {
		minutes = k.EstimatedMinutes
		if req.Type == model.EventPractice {
			minutes = maxInt(req.Total*2, 3)
		}
	}
	_ = s.repo.AddEvent(&model.LearningEvent{
		ID: genID("ev"), UserID: userID, KnowledgeID: k.ID, Type: req.Type,
		Correct: req.Correct, Total: req.Total, Source: req.Source, Minutes: minutes, CreatedAt: now,
	})
	v := s.stateView(c, *st, now)
	return &v, nil
}

// Practice 对知识点练习题判分并记录练习事件。
func (s *GrowthService) Practice(userID, knowledgeID string, answers []int) (*model.PracticeResult, error) {
	k, err := s.catalog.KnowledgeFull(knowledgeID)
	if err != nil {
		return nil, err
	}
	var qs []model.Question
	_ = json.Unmarshal(k.Questions, &qs)
	if len(qs) == 0 {
		return nil, apperr.BadRequest("该知识点暂无练习题")
	}
	res := &model.PracticeResult{Total: len(qs), Results: make([]bool, len(qs))}
	for i, q := range qs {
		if i < len(answers) && answers[i] == q.Answer {
			res.Correct++
			res.Results[i] = true
		}
	}
	view, err := s.RecordEvent(userID, model.LearningEventRequest{
		KnowledgeID: knowledgeID, Type: model.EventPractice, Correct: res.Correct, Total: res.Total, Source: "knowledge-practice",
	})
	if err != nil {
		return nil, err
	}
	res.State = *view
	return res, nil
}

// CompleteChapter 课程章节学习/测验结果映射到知识状态：课程结构与知识结构分离，经映射表传导。
func (s *GrowthService) CompleteChapter(userID string, req model.ChapterCompleteRequest) ([]model.KnowledgeStateView, error) {
	ids := s.catalog.ChapterKnowledge(req.CourseID, req.ChapterID)
	if len(ids) == 0 {
		return []model.KnowledgeStateView{}, nil
	}
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	out := make([]model.KnowledgeStateView, 0, len(ids))
	var skills []string
	for _, kid := range ids {
		ev := model.LearningEventRequest{KnowledgeID: kid, Type: model.EventStudy, Source: "course:" + req.CourseID}
		if req.Total > 0 {
			ev.Type, ev.Correct, ev.Total = model.EventPractice, req.Correct, req.Total
			ev.Source = "quiz:" + req.CourseID
		}
		v, err := s.applyEvent(c, userID, ev)
		if err != nil {
			continue
		}
		out = append(out, *v)
		skills = append(skills, c.KnowledgeSkills[kid]...)
	}
	s.recomputeSkills(c, userID, skills)
	return out, nil
}

func (s *GrowthService) knowledgeStateMap(userID string) map[string]model.StudentKnowledgeState {
	list, _ := s.repo.ListKnowledgeStates(userID)
	m := make(map[string]model.StudentKnowledgeState, len(list))
	for _, st := range list {
		m[st.KnowledgeID] = st
	}
	return m
}

// effectiveMasteryMap 返回全部知识点当前（含遗忘衰减）掌握度。
func (s *GrowthService) effectiveMasteryMap(userID string) map[string]float64 {
	now := s.now()
	out := map[string]float64{}
	for id, st := range s.knowledgeStateMap(userID) {
		out[id] = growth.EffectiveMastery(toEngineState(st), now)
	}
	return out
}

func (s *GrowthService) KnowledgeStates(userID string) ([]model.KnowledgeStateView, error) {
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	now := s.now()
	list, _ := s.repo.ListKnowledgeStates(userID)
	out := make([]model.KnowledgeStateView, 0, len(list))
	for _, st := range list {
		if _, ok := c.Knowledge[st.KnowledgeID]; ok {
			out = append(out, s.stateView(c, st, now))
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].LastStudiedAt > out[j].LastStudiedAt })
	return out, nil
}

// ---------- 技能状态 ----------

// RecomputeSkills 重新计算用户技能状态；skillIDs 为空时重算全部技能。
func (s *GrowthService) RecomputeSkills(userID string, skillIDs []string) error {
	c, err := s.catalog.Load()
	if err != nil {
		return err
	}
	if len(skillIDs) == 0 {
		for _, sk := range c.SkillList {
			skillIDs = append(skillIDs, sk.ID)
		}
	}
	s.recomputeSkills(c, userID, skillIDs)
	return nil
}

func (s *GrowthService) recomputeSkills(c *Catalog, userID string, skillIDs []string) int {
	if len(skillIDs) == 0 {
		return 0
	}
	mastery := s.effectiveMasteryMap(userID)
	evidence, _ := s.repo.ListEvidence(userID)
	projectScores, taskScores := map[string][]int{}, map[string][]int{}
	for _, e := range evidence {
		if e.SourceType == model.EvidenceOpportunity {
			taskScores[e.SkillID] = append(taskScores[e.SkillID], e.Score)
		} else {
			projectScores[e.SkillID] = append(projectScores[e.SkillID], e.Score)
		}
	}
	now := s.now()
	seen := map[string]bool{}
	updated := 0
	for _, sid := range skillIDs {
		if seen[sid] {
			continue
		}
		seen[sid] = true
		if _, ok := c.Skills[sid]; !ok {
			continue
		}
		in := growth.SkillInput{ProjectScores: projectScores[sid], TaskScores: taskScores[sid]}
		for _, link := range c.SkillKnowledge[sid] {
			in.Knowledge = append(in.Knowledge, growth.WeightedMastery{Mastery: mastery[link.KnowledgeID], Weight: link.Weight})
		}
		r := growth.EvaluateSkill(in)
		if err := s.repo.UpsertSkillState(&model.StudentSkillState{
			ID: genID("ss"), UserID: userID, SkillID: sid, Level: r.Level, Score: r.Score,
			KnowledgeScore: r.KnowledgeScore, ProjectScore: r.ProjectScore, TaskScore: r.TaskScore,
			EvidenceCount: r.EvidenceCount, UpdatedAt: now,
		}); err == nil {
			updated++
		}
	}
	return updated
}

// SkillLevels 返回用户各技能当前等级。
func (s *GrowthService) SkillLevels(userID string) map[string]int {
	list, _ := s.repo.ListSkillStates(userID)
	out := make(map[string]int, len(list))
	for _, st := range list {
		out[st.SkillID] = st.Level
	}
	return out
}

func skillStateView(sk model.Skill, st model.StudentSkillState) model.SkillStateView {
	return model.SkillStateView{
		Skill: sk, Level: st.Level, LevelName: growth.LevelNames[clampLevel(st.Level)], Score: st.Score,
		KnowledgeScore: st.KnowledgeScore, ProjectScore: st.ProjectScore, TaskScore: st.TaskScore, EvidenceCount: st.EvidenceCount,
	}
}

// SkillStates 返回用户全部技能状态（未开始的技能为 L0），按等级与得分排序。
func (s *GrowthService) SkillStates(userID string) ([]model.SkillStateView, error) {
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	return s.skillStatesWith(c, userID), nil
}

func (s *GrowthService) skillStatesWith(c *Catalog, userID string) []model.SkillStateView {
	list, _ := s.repo.ListSkillStates(userID)
	byID := map[string]model.StudentSkillState{}
	for _, st := range list {
		byID[st.SkillID] = st
	}
	out := make([]model.SkillStateView, 0, len(c.SkillList))
	for _, sk := range c.SkillList {
		out = append(out, skillStateView(sk, byID[sk.ID]))
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Level != out[j].Level {
			return out[i].Level > out[j].Level
		}
		return out[i].Score > out[j].Score
	})
	return out
}

// AddEvidence 写入能力证据并重算相关技能。
func (s *GrowthService) AddEvidence(userID string, skillIDs []string, sourceType, sourceID, title string, score int, detail string) []model.Evidence {
	now := s.now()
	var out []model.Evidence
	for _, sid := range skillIDs {
		e := model.Evidence{
			ID: genID("evd"), UserID: userID, SkillID: sid, SourceType: sourceType, SourceID: sourceID,
			Title: title, Score: score, Detail: detail, CreatedAt: now,
		}
		if err := s.repo.AddEvidence(&e); err == nil {
			out = append(out, e)
		}
	}
	_ = s.RecomputeSkills(userID, skillIDs)
	return out
}

// ---------- 差距分析 ----------

// Gap 分析用户相对目标岗位的能力差距；roleID 为空时使用当前职业目标。
func (s *GrowthService) Gap(userID, roleID string) (*model.GapAnalysis, error) {
	if roleID == "" {
		g := s.Goal(userID)
		if g == nil {
			return nil, nil
		}
		roleID = g.RoleID
	}
	career, role, err := s.catalog.RoleFull(roleID)
	if err != nil {
		return nil, err
	}
	return buildGap(career, role, s.SkillLevels(userID)), nil
}

func buildGap(career *model.Career, role *model.JobRole, levels map[string]int) *model.GapAnalysis {
	res := &model.GapAnalysis{Career: career, Role: role}
	var weighted, wsum float64
	for _, cap := range role.Capabilities {
		reqs := make([]growth.Requirement, 0, len(cap.Skills))
		names := map[string]string{}
		for _, cs := range cap.Skills {
			reqs = append(reqs, growth.Requirement{SkillID: cs.SkillID, Required: cs.RequiredLevel, Weight: cs.Weight})
			if cs.Skill != nil {
				names[cs.SkillID] = cs.Skill.Name
			}
		}
		fit := growth.Fit(reqs, levels)
		cr := model.CapabilityReadiness{ID: cap.ID, Name: cap.Name, Weight: maxInt(cap.Weight, 1), Readiness: fit.Score}
		for _, it := range fit.Items {
			item := model.SkillGapItem{
				SkillID: it.SkillID, SkillName: names[it.SkillID], CapabilityID: cap.ID, CapabilityName: cap.Name,
				Required: it.Required, Current: it.Current, Gap: it.Gap, Fit: it.Fit, Weight: it.Weight,
			}
			cr.Skills = append(cr.Skills, item)
			res.TotalCount++
			if it.Gap == 0 {
				res.MetCount++
			} else {
				item.Weight *= cr.Weight
				res.Gaps = append(res.Gaps, item)
			}
		}
		res.Capabilities = append(res.Capabilities, cr)
		weighted += float64(fit.Score * cr.Weight)
		wsum += float64(cr.Weight)
	}
	if wsum > 0 {
		res.Readiness = int(weighted/wsum + 0.5)
	}
	sort.SliceStable(res.Gaps, func(i, j int) bool {
		return res.Gaps[i].Gap*res.Gaps[i].Weight > res.Gaps[j].Gap*res.Gaps[j].Weight
	})
	res.Gaps = mergeGaps(res.Gaps)
	return res
}

// mergeGaps 同一技能出现在多个能力中时只保留差距最大的一条。
func mergeGaps(items []model.SkillGapItem) []model.SkillGapItem {
	seen := map[string]bool{}
	out := make([]model.SkillGapItem, 0, len(items))
	for _, it := range items {
		if seen[it.SkillID] {
			continue
		}
		seen[it.SkillID] = true
		out = append(out, it)
	}
	return out
}

// ---------- 知识推荐与学习计划 ----------

// Recommend 根据技能差距与知识图谱推荐下一步学习的知识点。
func (s *GrowthService) Recommend(userID string, gap *model.GapAnalysis, limit int) ([]model.KnowledgeRecommendation, error) {
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	return s.recommendWith(c, userID, gap, limit), nil
}

func (s *GrowthService) recommendWith(c *Catalog, userID string, gap *model.GapAnalysis, limit int) []model.KnowledgeRecommendation {
	mastery := s.effectiveMasteryMap(userID)
	var cands []growth.Candidate
	reasons := map[string]string{}
	add := func(skillID string, priority float64, reason string) {
		for _, link := range c.SkillKnowledge[skillID] {
			k := c.Knowledge[link.KnowledgeID]
			cands = append(cands, growth.Candidate{
				KnowledgeID: k.ID, SkillID: skillID, Priority: priority, Mastery: mastery[k.ID],
				Minutes: k.EstimatedMinutes, Difficulty: k.Difficulty,
			})
			if _, ok := reasons[k.ID]; !ok {
				reasons[k.ID] = reason
			}
		}
	}
	if gap != nil && len(gap.Gaps) > 0 {
		for _, g := range gap.Gaps {
			add(g.SkillID, float64(g.Gap*maxInt(g.Weight, 1)),
				fmt.Sprintf("补齐「%s」L%d → L%d", g.SkillName, g.Current, g.Required))
		}
	} else {
		for _, sk := range c.SkillList {
			add(sk.ID, 1, "巩固「"+sk.Name+"」")
		}
	}
	recs := growth.Recommend(cands, c.Prereqs, mastery, limit)
	out := make([]model.KnowledgeRecommendation, 0, len(recs))
	for _, r := range recs {
		item := model.KnowledgeRecommendation{
			Knowledge: c.Brief(r.KnowledgeID), SkillID: r.SkillID, SkillName: c.SkillName(r.SkillID),
			Mastery: r.Mastery, Ready: r.Ready, Reason: reasons[r.KnowledgeID],
		}
		if item.Reason == "" {
			item.Reason = "后续知识的前置基础"
		}
		for _, p := range r.MissingPrereqs {
			item.MissingPrereqs = append(item.MissingPrereqs, c.Brief(p))
		}
		out = append(out, item)
	}
	return out
}

// BuildPlan 将推荐知识按周投入拆分为学习计划，并在每个技能达到 L2 后插入对应项目实践。
func (s *GrowthService) BuildPlan(userID string, goal *model.UserGoal, gap *model.GapAnalysis, recs []model.KnowledgeRecommendation) *model.LearningPlan {
	weekly := 600
	if goal != nil {
		weekly = goal.WeeklyHours * 60
	}
	items := make([]growth.PlanItem, 0, len(recs))
	byID := map[string]model.KnowledgeRecommendation{}
	for _, r := range recs {
		items = append(items, growth.PlanItem{KnowledgeID: r.Knowledge.ID, Minutes: r.Knowledge.EstimatedMinutes})
		byID[r.Knowledge.ID] = r
	}
	plan := &model.LearningPlan{}
	for i, wk := range growth.PlanWeeks(items, weekly) {
		pw := model.PlanWeek{Week: i + 1}
		skillCount := map[string]int{}
		for _, it := range wk {
			r := byID[it.KnowledgeID]
			pw.Items = append(pw.Items, model.PlanItem{
				Type: "knowledge", RefID: it.KnowledgeID, Title: r.Knowledge.Name, Minutes: it.Minutes, Reason: r.Reason,
			})
			plan.TotalMinutes += it.Minutes
			skillCount[r.SkillName]++
		}
		pw.Focus = topKey(skillCount)
		plan.Weeks = append(plan.Weeks, pw)
	}
	if task := s.bestPracticeTask(userID, gap); task != nil {
		minutes := 180
		item := model.PlanItem{Type: "project", RefID: task.ProjectID, SubID: task.TaskID, Title: task.ProjectTitle + " · " + task.TaskTitle, Minutes: minutes, Reason: "以项目实践把知识转化为技能证据"}
		if len(plan.Weeks) == 0 {
			plan.Weeks = append(plan.Weeks, model.PlanWeek{Week: 1, Focus: "项目实践"})
		}
		last := &plan.Weeks[len(plan.Weeks)-1]
		last.Items = append(last.Items, item)
		plan.TotalMinutes += minutes
	}
	role := "当前方向"
	if gap != nil && gap.Role != nil {
		role = gap.Role.Name
	}
	plan.Summary = fmt.Sprintf("围绕「%s」规划 %d 周、共 %d 项学习任务，预计投入 %.1f 小时。", role, len(plan.Weeks), countItems(plan), float64(plan.TotalMinutes)/60)
	return plan
}

func countItems(p *model.LearningPlan) int {
	n := 0
	for _, w := range p.Weeks {
		n += len(w.Items)
	}
	return n
}

func topKey(m map[string]int) string {
	best, n := "", -1
	for k, v := range m {
		if v > n || (v == n && k < best) {
			best, n = k, v
		}
	}
	return best
}

// bestPracticeTask 推荐与差距技能最相关、且当前能力接近要求的项目任务。
func (s *GrowthService) bestPracticeTask(userID string, gap *model.GapAnalysis) *model.ProjectTaskRef {
	refs := s.PracticeTasks(userID, gap, 1)
	if len(refs) == 0 {
		return nil
	}
	return &refs[0]
}

// PracticeTasks 推荐项目任务：覆盖差距技能越多、匹配度越接近可完成区间越优先。
func (s *GrowthService) PracticeTasks(userID string, gap *model.GapAnalysis, limit int) []model.ProjectTaskRef {
	projects, _ := s.projects.List("published")
	c, _ := s.catalog.Load()
	levels := s.SkillLevels(userID)
	gapSkills := map[string]bool{}
	if gap != nil {
		for _, g := range gap.Gaps {
			gapSkills[g.SkillID] = true
		}
	}
	done := map[string]bool{}
	subs, _ := s.projects.ListSubmissions(userID, "", 0)
	for _, sub := range subs {
		if sub.Score >= 60 {
			done[sub.TaskID] = true
		}
	}
	type scored struct {
		ref   model.ProjectTaskRef
		score float64
	}
	var all []scored
	for _, p := range projects {
		for _, t := range p.Tasks {
			if done[t.ID] {
				continue
			}
			reqs := parseRequirements(t.Requirements)
			fit := growth.Fit(toEngineReqs(reqs), levels)
			relevance := 0
			for _, r := range reqs {
				if gapSkills[r.SkillID] {
					relevance++
				}
			}
			if len(gapSkills) > 0 && relevance == 0 {
				continue
			}
			// 最适合练习的是"跳一跳够得着"的任务：匹配度 50-90 之间优先
			stretch := 1 - absF(float64(fit.Score)-75)/75
			all = append(all, scored{
				ref:   model.ProjectTaskRef{ProjectID: p.ID, ProjectTitle: p.Title, TaskID: t.ID, TaskTitle: t.Title, Match: matchView(fit, c)},
				score: float64(relevance)*2 + stretch,
			})
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].score > all[j].score })
	out := make([]model.ProjectTaskRef, 0, limit)
	for i := 0; i < len(all) && (limit <= 0 || i < limit); i++ {
		out = append(out, all[i].ref)
	}
	return out
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func toEngineReqs(reqs []model.SkillRequirement) []growth.Requirement {
	out := make([]growth.Requirement, 0, len(reqs))
	for _, r := range reqs {
		out = append(out, growth.Requirement{SkillID: r.SkillID, Required: r.Level, Weight: 1})
	}
	return out
}

func matchView(fit growth.FitResult, c *Catalog) model.MatchView {
	mv := model.MatchView{Score: fit.Score, Qualified: fit.Qualified, MetCount: fit.MetCount}
	for _, it := range fit.Items {
		item := model.SkillGapItem{SkillID: it.SkillID, Required: it.Required, Current: it.Current, Gap: it.Gap, Fit: it.Fit, Weight: it.Weight}
		if c != nil {
			item.SkillName = c.SkillName(it.SkillID)
		}
		mv.Items = append(mv.Items, item)
	}
	return mv
}

// ---------- 知识图谱 ----------

// Graph 返回带个人掌握度的知识图谱；skillID 不为空时只返回该技能相关知识及其前置。
func (s *GrowthService) Graph(userID, skillID string) (*model.KnowledgeGraph, error) {
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	include := map[string]bool{}
	if skillID != "" {
		var walk func(id string)
		walk = func(id string) {
			if include[id] {
				return
			}
			include[id] = true
			for _, p := range c.Prereqs[id] {
				walk(p)
			}
		}
		for _, l := range c.SkillKnowledge[skillID] {
			walk(l.KnowledgeID)
		}
	}
	states := s.knowledgeStateMap(userID)
	depth := growth.Depths(c.Prereqs)
	now := s.now()
	g := &model.KnowledgeGraph{}
	for _, k := range c.KnowledgeList {
		if skillID != "" && !include[k.ID] {
			continue
		}
		st, ok := states[k.ID]
		node := model.GraphNode{
			KnowledgeBrief: knowledgeBrief(k), Summary: k.Summary, SkillIDs: nonNil(c.KnowledgeSkills[k.ID]),
			Depth: depth[k.ID], Status: model.KnowledgeStatusNew,
		}
		if ok {
			node.Mastery = growth.EffectiveMastery(toEngineState(st), now)
			node.Status = st.Status
		}
		g.Nodes = append(g.Nodes, node)
	}
	for _, r := range c.Relations {
		if skillID != "" && (!include[r.FromID] || !include[r.ToID]) {
			continue
		}
		g.Edges = append(g.Edges, r)
	}
	g.Nodes, g.Edges = nonNil(g.Nodes), nonNil(g.Edges)
	return g, nil
}

// KnowledgeDetail 知识点详情：正文、练习、关系、技能、课程与知识资产。
func (s *GrowthService) KnowledgeDetail(userID, id string, resources []model.KnowledgeResource) (*model.KnowledgeDetail, error) {
	k, err := s.catalog.KnowledgeFull(id)
	if err != nil {
		return nil, err
	}
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	var qs []model.Question
	_ = json.Unmarshal(k.Questions, &qs)
	detail := &model.KnowledgeDetail{Knowledge: *k, Resources: nonNil(resources), Chapters: nonNil(s.catalog.ChapterRefs(id))}
	detail.Knowledge.Questions = nil
	for _, q := range qs {
		detail.Questions = append(detail.Questions, model.PracticeQuestion{Text: q.Text, Options: q.Options})
	}
	detail.Questions = nonNil(detail.Questions)
	now := s.now()
	states := s.knowledgeStateMap(userID)
	if st, ok := states[id]; ok {
		detail.State = s.stateView(c, st, now)
	} else {
		detail.State = model.KnowledgeStateView{Knowledge: c.Brief(id), Status: model.KnowledgeStatusNew}
	}
	for _, sid := range c.KnowledgeSkills[id] {
		if sk, ok := c.Skills[sid]; ok {
			detail.Skills = append(detail.Skills, sk)
		}
	}
	detail.Skills = nonNil(detail.Skills)
	for _, r := range c.Relations {
		var link *model.KnowledgeLink
		switch id {
		case r.ToID:
			link = &model.KnowledgeLink{Type: r.Type, Direction: "in", Knowledge: c.Brief(r.FromID)}
		case r.FromID:
			link = &model.KnowledgeLink{Type: r.Type, Direction: "out", Knowledge: c.Brief(r.ToID)}
		}
		if link != nil {
			if st, ok := states[link.Knowledge.ID]; ok {
				link.Mastery = growth.EffectiveMastery(toEngineState(st), now)
			}
			detail.Relations = append(detail.Relations, *link)
		}
	}
	detail.Relations = nonNil(detail.Relations)
	return detail, nil
}

// ---------- 成长飞轮与活跃度 ----------

func (s *GrowthService) Flywheel(userID string, contributions int) model.FlywheelStats {
	f := model.FlywheelStats{LearningEvents: s.repo.CountEvents(userID), Contributions: contributions}
	events, _ := s.repo.ListEvents(userID, 0)
	for _, e := range events {
		f.LearningMinutes += e.Minutes
	}
	for _, st := range s.knowledgeStateMap(userID) {
		if st.Status == model.KnowledgeStatusMastered {
			f.MasteredCount++
		}
	}
	subs, _ := s.projects.ListSubmissions(userID, "", 0)
	projects := map[string]bool{}
	for _, sub := range subs {
		projects[sub.ProjectID] = true
	}
	f.ProjectsPracticed, f.Submissions = len(projects), len(subs)
	ev, _ := s.repo.ListEvidence(userID)
	f.EvidenceCount = len(ev)
	apps, _ := s.market.ListApplications("", userID)
	for _, a := range apps {
		if a.Status == model.ApplicationCompleted {
			f.TasksCompleted++
		}
	}
	income, _ := s.repo.ListIncome(userID)
	for _, i := range income {
		f.Income += i.Amount
	}
	states, _ := s.repo.ListSkillStates(userID)
	total, n := 0, 0
	for _, st := range states {
		if st.Level > 0 {
			total += st.Level
			n++
		}
	}
	if n > 0 {
		f.AvgSkillLevel = float64(int(float64(total)/float64(n)*10+0.5)) / 10
	}
	return f
}

// Activity 最近 days 天的学习活跃度。
func (s *GrowthService) Activity(userID string, days int) []model.ActivityDay {
	start := time.Now().AddDate(0, 0, -(days - 1))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	events, _ := s.repo.ListEvents(userID, start.UnixMilli())
	byDay := map[string]*model.ActivityDay{}
	out := make([]model.ActivityDay, days)
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i).Format("2006-01-02")
		out[i] = model.ActivityDay{Date: d}
		byDay[d] = &out[i]
	}
	for _, e := range events {
		d := time.UnixMilli(e.CreatedAt).Format("2006-01-02")
		if day, ok := byDay[d]; ok {
			day.Events++
			day.Minutes += e.Minutes
		}
	}
	return out
}

// DueReviews 到期需要复习的知识点。
func (s *GrowthService) DueReviews(userID string, limit int) []model.KnowledgeStateView {
	views, _ := s.KnowledgeStates(userID)
	out := []model.KnowledgeStateView{}
	for _, v := range views {
		if v.Due {
			out = append(out, v)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].NextReviewAt < out[j].NextReviewAt })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// Overview 成长中心首页数据。
func (s *GrowthService) Overview(userID string, opportunities []model.OpportunityView, contributions int) (*model.GrowthOverview, error) {
	goal := s.Goal(userID)
	gap, err := s.Gap(userID, "")
	if err != nil {
		gap = nil
	}
	recs, err := s.Recommend(userID, gap, 6)
	if err != nil {
		return nil, err
	}
	ov := &model.GrowthOverview{
		Goal: goal, Gap: gap, Recommendations: recs,
		DueReviews: s.DueReviews(userID, 6), Flywheel: s.Flywheel(userID, contributions),
		Activity: s.Activity(userID, 28), NextProjects: nonNil(s.PracticeTasks(userID, gap, 3)),
		Opportunities: nonNil(opportunities),
	}
	if run, err := s.kernel.Latest(userID, "learning"); err == nil {
		ov.LatestPlan = run
	}
	return ov, nil
}
