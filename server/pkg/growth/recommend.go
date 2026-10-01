package growth

import "sort"

// Candidate 待推荐的知识点。
type Candidate struct {
	KnowledgeID string
	SkillID     string
	Priority    float64 // 来自技能差距的优先级，越大越优先
	Mastery     float64
	Minutes     int
	Difficulty  int
}

// Recommendation 知识推荐结果。
type Recommendation struct {
	KnowledgeID    string
	SkillID        string
	Mastery        float64
	Ready          bool
	MissingPrereqs []string
	Minutes        int
}

// PrereqReadyThreshold 前置知识达到该掌握度即视为可以进入后续知识。
const PrereqReadyThreshold = 0.6

// MasteredThreshold 达到该掌握度的知识不再推荐学习。
const MasteredThreshold = 0.8

// Recommend 基于知识图谱前置关系推荐下一步学习的知识点。
//
// prereqs[k] 为 k 的前置知识列表，mastery 为全部知识的当前掌握度。
// 未满足的前置知识会被自动补入候选，排序规则：已就绪优先 → 技能差距优先级 → 图谱深度（先基础后进阶）→ 难度。
func Recommend(cands []Candidate, prereqs map[string][]string, mastery map[string]float64, limit int) []Recommendation {
	byID := map[string]Candidate{}
	for _, c := range cands {
		if old, ok := byID[c.KnowledgeID]; !ok || c.Priority > old.Priority {
			byID[c.KnowledgeID] = c
		}
	}
	// 递归补入未掌握的前置知识，继承后继的优先级
	var pull func(id string, priority float64, skill string, seen map[string]bool)
	pull = func(id string, priority float64, skill string, seen map[string]bool) {
		for _, p := range prereqs[id] {
			if seen[p] || mastery[p] >= PrereqReadyThreshold {
				continue
			}
			seen[p] = true
			if old, ok := byID[p]; !ok || old.Priority < priority {
				c := old
				c.KnowledgeID, c.Priority, c.Mastery = p, priority, mastery[p]
				if c.SkillID == "" {
					c.SkillID = skill
				}
				byID[p] = c
			}
			pull(p, priority, skill, seen)
		}
	}
	for _, c := range cands {
		pull(c.KnowledgeID, c.Priority, c.SkillID, map[string]bool{c.KnowledgeID: true})
	}

	depth := Depths(prereqs)
	out := make([]Recommendation, 0, len(byID))
	for id, c := range byID {
		m := mastery[id]
		if m >= MasteredThreshold {
			continue
		}
		rec := Recommendation{KnowledgeID: id, SkillID: c.SkillID, Mastery: m, Ready: true, Minutes: c.Minutes}
		for _, p := range prereqs[id] {
			if mastery[p] < PrereqReadyThreshold {
				rec.Ready = false
				rec.MissingPrereqs = append(rec.MissingPrereqs, p)
			}
		}
		out = append(out, rec)
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Ready != b.Ready {
			return a.Ready
		}
		pa, pb := byID[a.KnowledgeID].Priority, byID[b.KnowledgeID].Priority
		if pa != pb {
			return pa > pb
		}
		if depth[a.KnowledgeID] != depth[b.KnowledgeID] {
			return depth[a.KnowledgeID] < depth[b.KnowledgeID]
		}
		if byID[a.KnowledgeID].Difficulty != byID[b.KnowledgeID].Difficulty {
			return byID[a.KnowledgeID].Difficulty < byID[b.KnowledgeID].Difficulty
		}
		return a.KnowledgeID < b.KnowledgeID
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// Depths 计算每个知识点在前置关系图中的深度（无前置为 0），环路按已访问截断。
func Depths(prereqs map[string][]string) map[string]int {
	memo := map[string]int{}
	var visit func(id string, stack map[string]bool) int
	visit = func(id string, stack map[string]bool) int {
		if d, ok := memo[id]; ok {
			return d
		}
		if stack[id] {
			return 0
		}
		stack[id] = true
		d := 0
		for _, p := range prereqs[id] {
			if v := visit(p, stack) + 1; v > d {
				d = v
			}
		}
		delete(stack, id)
		memo[id] = d
		return d
	}
	for id := range prereqs {
		visit(id, map[string]bool{})
	}
	return memo
}

// PlanItem 学习计划中的一项。
type PlanItem struct {
	KnowledgeID string
	Minutes     int
}

// PlanWeeks 按每周可投入时间把推荐知识切分为周计划，单项超出周预算时独占一周。
func PlanWeeks(items []PlanItem, weeklyMinutes int) [][]PlanItem {
	if weeklyMinutes <= 0 {
		weeklyMinutes = 300
	}
	var weeks [][]PlanItem
	var cur []PlanItem
	used := 0
	for _, it := range items {
		m := it.Minutes
		if m <= 0 {
			m = 30
		}
		if used+m > weeklyMinutes && len(cur) > 0 {
			weeks = append(weeks, cur)
			cur, used = nil, 0
		}
		cur = append(cur, it)
		used += m
	}
	if len(cur) > 0 {
		weeks = append(weeks, cur)
	}
	return weeks
}
