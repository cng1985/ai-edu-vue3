package growth

import "testing"

func TestApplyStudyOnlyCapsMastery(t *testing.T) {
	s := KnowledgeState{}
	for i := 0; i < 20; i++ {
		s = Apply(s, Event{Type: "study", At: 1})
	}
	if s.Mastery > 0.3 {
		t.Fatalf("study only should cap mastery at 0.3, got %v", s.Mastery)
	}
	if s.Status != StatusLearning {
		t.Fatalf("expected learning, got %s", s.Status)
	}
}

func TestApplyPracticeRaisesMasteryAndConfidence(t *testing.T) {
	s := Apply(KnowledgeState{}, Event{Type: "study", At: 1})
	s = Apply(s, Event{Type: "practice", Correct: 9, Total: 10, At: 2})
	s = Apply(s, Event{Type: "practice", Correct: 10, Total: 10, At: 3})
	if s.Mastery < 0.8 {
		t.Fatalf("expected mastery >= 0.8, got %v", s.Mastery)
	}
	if s.Status != StatusMastered {
		t.Fatalf("expected mastered, got %s (confidence %v)", s.Status, s.Confidence)
	}
	if s.NextReviewAt <= s.LastStudiedAt {
		t.Fatalf("next review should be scheduled after last study")
	}

	low := Apply(KnowledgeState{}, Event{Type: "practice", Correct: 1, Total: 10, At: 1})
	if low.Mastery >= 0.3 {
		t.Fatalf("poor practice should keep mastery low, got %v", low.Mastery)
	}
	if ReviewIntervalDays(low.Mastery) != 1 {
		t.Fatalf("low mastery should be reviewed next day")
	}
}

func TestEffectiveMasteryDecays(t *testing.T) {
	s := KnowledgeState{Mastery: 0.9, LastStudiedAt: 1}
	later := EffectiveMastery(s, 1+60*dayMs)
	if later >= 0.9 || later < 0.54 {
		t.Fatalf("unexpected decayed mastery %v", later)
	}
	if EffectiveMastery(s, 1) != 0.9 {
		t.Fatalf("no decay without elapsed time")
	}
}

func TestEvaluateSkillKnowledgeIsNotSkill(t *testing.T) {
	r := EvaluateSkill(SkillInput{Knowledge: []WeightedMastery{{Mastery: 1, Weight: 1}}})
	if r.Level != 2 {
		t.Fatalf("knowledge only should cap at L2, got L%d", r.Level)
	}
	r = EvaluateSkill(SkillInput{Knowledge: []WeightedMastery{{Mastery: 0.85, Weight: 1}}, ProjectScores: []int{78}})
	if r.Level != 3 {
		t.Fatalf("project evidence should unlock L3, got L%d (score %v)", r.Level, r.Score)
	}
	r = EvaluateSkill(SkillInput{
		Knowledge:     []WeightedMastery{{Mastery: 0.95, Weight: 1}},
		ProjectScores: []int{92, 88},
		TaskScores:    []int{90, 95},
	})
	if r.Level != 5 {
		t.Fatalf("excellent evidence should reach L5, got L%d (score %v)", r.Level, r.Score)
	}
	r = EvaluateSkill(SkillInput{Knowledge: []WeightedMastery{{Mastery: 0.3, Weight: 1}}, ProjectScores: []int{90}, TaskScores: []int{96}})
	if r.Level != 2 {
		t.Fatalf("weak knowledge should cap evidence-backed skill at L2, got L%d", r.Level)
	}
	if EvaluateSkill(SkillInput{}).Level != 0 {
		t.Fatalf("no data should be L0")
	}
}

func TestFitAndRankGaps(t *testing.T) {
	reqs := []Requirement{{SkillID: "spring", Required: 3, Weight: 2}, {SkillID: "mysql", Required: 3, Weight: 1}, {SkillID: "redis", Required: 2, Weight: 1}}
	res := Fit(reqs, map[string]int{"spring": 3, "mysql": 1, "redis": 2})
	if res.Qualified {
		t.Fatalf("should not be qualified")
	}
	if res.MetCount != 2 {
		t.Fatalf("expected 2 met, got %d", res.MetCount)
	}
	// spring 1*2 + mysql (1/3)/2*1 + redis 1*1 = 3.1667 / 4 ≈ 79
	if res.Score != 79 {
		t.Fatalf("unexpected score %d", res.Score)
	}
	gaps := RankGaps(res.Items)
	if len(gaps) != 1 || gaps[0].SkillID != "mysql" || gaps[0].Gap != 2 {
		t.Fatalf("unexpected gaps %+v", gaps)
	}
	full := Fit(reqs, map[string]int{"spring": 4, "mysql": 3, "redis": 5})
	if !full.Qualified || full.Score != 100 {
		t.Fatalf("expected qualified 100, got %+v", full)
	}
}

func TestRecommendRespectsPrerequisites(t *testing.T) {
	prereqs := map[string][]string{
		"aqs":      {"cas"},
		"cas":      {"jmm"},
		"volatile": {"jmm"},
	}
	cands := []Candidate{
		{KnowledgeID: "aqs", SkillID: "concurrency", Priority: 2},
		{KnowledgeID: "volatile", SkillID: "concurrency", Priority: 2},
	}
	recs := Recommend(cands, prereqs, map[string]float64{}, 0)
	if len(recs) != 4 {
		t.Fatalf("expected prerequisites to be pulled in, got %+v", recs)
	}
	if recs[0].KnowledgeID != "jmm" || !recs[0].Ready {
		t.Fatalf("jmm should be first and ready, got %+v", recs[0])
	}
	for _, r := range recs {
		if r.KnowledgeID == "aqs" && (r.Ready || len(r.MissingPrereqs) != 1) {
			t.Fatalf("aqs should be blocked by cas: %+v", r)
		}
	}

	recs = Recommend(cands, prereqs, map[string]float64{"jmm": 0.9, "cas": 0.7, "volatile": 0.85}, 0)
	if len(recs) != 1 || recs[0].KnowledgeID != "aqs" || !recs[0].Ready {
		t.Fatalf("expected only aqs ready, got %+v", recs)
	}
}

func TestPlanWeeks(t *testing.T) {
	weeks := PlanWeeks([]PlanItem{{"a", 60}, {"b", 60}, {"c", 200}, {"d", 30}}, 150)
	if len(weeks) != 3 || len(weeks[0]) != 2 || len(weeks[1]) != 1 || len(weeks[2]) != 1 {
		t.Fatalf("unexpected weeks %+v", weeks)
	}
}
