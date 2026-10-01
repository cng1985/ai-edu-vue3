package growth

import "sort"

// Requirement 一项技能要求。
type Requirement struct {
	SkillID  string
	Required int
	Weight   int
}

// RequirementFit 单项技能要求的满足情况。
type RequirementFit struct {
	SkillID  string  `json:"skillId"`
	Required int     `json:"required"`
	Current  int     `json:"current"`
	Gap      int     `json:"gap"`
	Fit      float64 `json:"fit"`
	Weight   int     `json:"weight"`
}

// FitResult 一组要求的整体满足度。
type FitResult struct {
	Score     int              `json:"score"` // 0-100
	Qualified bool             `json:"qualified"`
	MetCount  int              `json:"metCount"`
	Items     []RequirementFit `json:"items"`
}

// Fit 计算当前技能等级对一组要求的满足度，用于岗位差距分析与任务匹配。
//
// 单项满足度 = min(当前/要求, 1)，差 2 级及以上额外折半，体现"差一点可以补，差太多难以胜任"。
func Fit(reqs []Requirement, levels map[string]int) FitResult {
	res := FitResult{Qualified: true}
	var sum, wsum float64
	for _, r := range reqs {
		cur := levels[r.SkillID]
		item := RequirementFit{SkillID: r.SkillID, Required: r.Required, Current: cur, Weight: r.Weight}
		if r.Required <= 0 {
			item.Fit = 1
		} else {
			item.Fit = clamp(float64(cur)/float64(r.Required), 0, 1)
			if r.Required-cur >= 2 {
				item.Fit /= 2
			}
		}
		if cur < r.Required {
			item.Gap = r.Required - cur
			res.Qualified = false
		} else {
			res.MetCount++
		}
		item.Fit = round3(item.Fit)
		w := float64(r.Weight)
		if w <= 0 {
			w = 1
		}
		sum += item.Fit * w
		wsum += w
		res.Items = append(res.Items, item)
	}
	if wsum > 0 {
		res.Score = int(sum/wsum*100 + 0.5)
	} else {
		res.Score = 100
	}
	return res
}

// RankGaps 按"差距 × 权重"从大到小排列未达标的技能。
func RankGaps(items []RequirementFit) []RequirementFit {
	out := make([]RequirementFit, 0, len(items))
	for _, it := range items {
		if it.Gap > 0 {
			out = append(out, it)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		wi, wj := out[i].Gap*maxInt(out[i].Weight, 1), out[j].Gap*maxInt(out[j].Weight, 1)
		if wi != wj {
			return wi > wj
		}
		return out[i].SkillID < out[j].SkillID
	})
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
