// Package growth 是个人成长计算引擎：知识掌握度、技能等级、能力差距、人才匹配与知识推荐。
// 全部为无副作用的纯函数，持久化由 service 层负责。
package growth

import (
	"math"
	"time"
)

const dayMs = int64(24 * time.Hour / time.Millisecond)

// KnowledgeState 知识掌握状态的计算输入/输出。
type KnowledgeState struct {
	StudyCount    int
	PracticeCount int
	CorrectCount  int
	Accuracy      float64
	Mastery       float64
	Confidence    float64
	Status        string
	LastStudiedAt int64
	NextReviewAt  int64
}

// Event 一次学习行为。Practice 事件的 Total 为题目数，Correct 为答对数。
type Event struct {
	Type    string // study | practice | review
	Correct int
	Total   int
	At      int64
}

const (
	StatusNew      = "new"
	StatusLearning = "learning"
	StatusMastered = "mastered"
)

// Apply 根据学习事件更新知识状态。
//
// 掌握度由两部分组成：学习投入（学习与练习次数，边际递减，最多贡献 0.3）与练习表现（贝叶斯平滑正确率 × 练习可信度）。
// 仅学习不练习时掌握度上限为 0.3，体现"看过 ≠ 学会"。
func Apply(s KnowledgeState, e Event) KnowledgeState {
	switch e.Type {
	case "practice":
		if e.Total > 0 {
			s.PracticeCount += e.Total
			s.CorrectCount += clampInt(e.Correct, 0, e.Total)
		}
	case "review":
		s.StudyCount++
		if e.Total > 0 {
			s.PracticeCount += e.Total
			s.CorrectCount += clampInt(e.Correct, 0, e.Total)
		}
	default:
		s.StudyCount++
	}

	if s.PracticeCount > 0 {
		s.Accuracy = float64(s.CorrectCount) / float64(s.PracticeCount)
	}
	engagement := 1 - math.Exp(-(float64(s.StudyCount)+float64(s.PracticeCount)/5)/2)
	smoothedAcc := (float64(s.CorrectCount) + 1) / (float64(s.PracticeCount) + 2)
	practiceTrust := float64(s.PracticeCount) / (float64(s.PracticeCount) + 2)
	s.Mastery = round3(clamp(0.3*engagement+0.7*smoothedAcc*practiceTrust, 0, 1))

	evidence := float64(s.StudyCount + s.PracticeCount)
	s.Confidence = round3(1 - 1/(1+0.3*evidence))

	s.Status = StatusOf(s.Mastery, s.Confidence)
	s.LastStudiedAt = e.At
	s.NextReviewAt = e.At + ReviewIntervalDays(s.Mastery)*dayMs
	return s
}

// StatusOf 由掌握度与置信度判定学习状态。
func StatusOf(mastery, confidence float64) string {
	switch {
	case mastery >= 0.8 && confidence >= 0.6:
		return StatusMastered
	case mastery > 0:
		return StatusLearning
	default:
		return StatusNew
	}
}

// ReviewIntervalDays 按掌握度给出间隔复习天数（1/3/7/15/30）。
func ReviewIntervalDays(mastery float64) int64 {
	switch {
	case mastery < 0.4:
		return 1
	case mastery < 0.6:
		return 3
	case mastery < 0.8:
		return 7
	case mastery < 0.9:
		return 15
	default:
		return 30
	}
}

// EffectiveMastery 考虑遗忘衰减后的当前掌握度：掌握越牢固，记忆稳定期越长；衰减下限为原值的 60%。
func EffectiveMastery(s KnowledgeState, now int64) float64 {
	if s.Mastery <= 0 || s.LastStudiedAt == 0 {
		return s.Mastery
	}
	days := float64(now-s.LastStudiedAt) / float64(dayMs)
	if days <= 0 {
		return s.Mastery
	}
	stability := 3 + 27*s.Mastery
	retention := 0.6 + 0.4*math.Exp(-days/stability)
	return round3(s.Mastery * retention)
}

// IsDue 是否到期需要复习。
func IsDue(s KnowledgeState, now int64) bool {
	return s.NextReviewAt > 0 && s.NextReviewAt <= now
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func round3(v float64) float64 { return math.Round(v*1000) / 1000 }
