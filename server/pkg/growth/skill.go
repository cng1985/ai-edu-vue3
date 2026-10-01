package growth

import "sort"

// LevelNames 技能等级定义。
var LevelNames = []string{
	"未掌握",
	"了解",
	"可以辅助完成",
	"可以独立完成",
	"可以解决复杂问题",
	"可以设计和指导他人",
}

// MaxLevel 最高技能等级。
const MaxLevel = 5

// WeightedMastery 技能关联的一个知识点及其权重。
type WeightedMastery struct {
	Mastery float64
	Weight  int
}

// SkillInput 技能状态计算输入：知识状态 + 项目证据 + 任务表现。
type SkillInput struct {
	Knowledge     []WeightedMastery
	ProjectScores []int // 0-100
	TaskScores    []int // 0-100
}

// SkillResult 技能状态计算结果，各分项为 0-1。
type SkillResult struct {
	Level          int
	Score          float64
	KnowledgeScore float64
	ProjectScore   float64
	TaskScore      float64
	EvidenceCount  int
}

// EvaluateSkill 计算技能等级。
//
// 知识 ≠ 技能：只有知识没有项目证据时最高 L2；L3 需要项目证据；L4 需要高质量项目或真实任务；
// L5 需要多项证据且真实任务表现优秀。反过来，证据也不能替代知识：L3/L4/L5 分别要求知识得分不低于 0.4/0.6/0.75，
// 知识薄弱但有证据时最高 L2。
func EvaluateSkill(in SkillInput) SkillResult {
	r := SkillResult{
		KnowledgeScore: weightedAvg(in.Knowledge),
		ProjectScore:   topAvg(in.ProjectScores, 3),
		TaskScore:      topAvg(in.TaskScores, 3),
		EvidenceCount:  len(in.ProjectScores) + len(in.TaskScores),
	}

	switch {
	case r.TaskScore > 0 && r.ProjectScore > 0:
		r.Score = 0.4*r.KnowledgeScore + 0.35*r.ProjectScore + 0.25*r.TaskScore
	case r.TaskScore > 0:
		r.Score = 0.5*r.KnowledgeScore + 0.5*r.TaskScore
	case r.ProjectScore > 0:
		r.Score = 0.55*r.KnowledgeScore + 0.45*r.ProjectScore
	default:
		r.Score = r.KnowledgeScore * 0.6
	}
	r.Score = round3(r.Score)

	best := maxF(r.ProjectScore, r.TaskScore)
	k := r.KnowledgeScore
	switch {
	case k >= 0.75 && r.Score >= 0.85 && r.TaskScore >= 0.85 && r.EvidenceCount >= 2:
		r.Level = 5
	case k >= 0.6 && r.Score >= 0.72 && (r.ProjectScore >= 0.8 || r.TaskScore >= 0.75):
		r.Level = 4
	case k >= 0.4 && r.Score >= 0.55 && best >= 0.6:
		r.Level = 3
	case k >= 0.5 || best >= 0.6:
		r.Level = 2
	case r.KnowledgeScore >= 0.15 || best > 0:
		r.Level = 1
	default:
		r.Level = 0
	}
	return r
}

func weightedAvg(items []WeightedMastery) float64 {
	var sum, w float64
	for _, it := range items {
		weight := float64(it.Weight)
		if weight <= 0 {
			weight = 1
		}
		sum += it.Mastery * weight
		w += weight
	}
	if w == 0 {
		return 0
	}
	return round3(sum / w)
}

// topAvg 取最好的 n 项证据的平均分，归一化到 0-1。
func topAvg(scores []int, n int) float64 {
	if len(scores) == 0 {
		return 0
	}
	sorted := append([]int(nil), scores...)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	total := 0
	for _, s := range sorted {
		total += clampInt(s, 0, 100)
	}
	return round3(float64(total) / float64(len(sorted)) / 100)
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
