package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type growthFixture struct {
	db     *gorm.DB
	growth *GrowthService
	market *MarketService
	gr     *repository.GrowthRepo
}

func newGrowthFixture(t *testing.T) *growthFixture {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Course{}, &model.Chapter{},
		&model.Career{}, &model.JobRole{}, &model.RoleCapability{}, &model.CapabilitySkill{},
		&model.Skill{}, &model.SkillKnowledge{}, &model.KnowledgePoint{}, &model.KnowledgeRelation{}, &model.ChapterKnowledge{},
		&model.StudentKnowledgeState{}, &model.LearningEvent{}, &model.StudentSkillState{}, &model.UserGoal{},
		&model.Project{}, &model.ProjectTask{}, &model.TaskSubmission{}, &model.Evidence{},
		&model.Opportunity{}, &model.OpportunityApplication{}, &model.IncomeRecord{}, &model.PipelineRun{},
	); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
	now := time.Now().UnixMilli()
	q := datatypes.JSON(`[{"text":"q1","options":["a","b"],"answer":1},{"text":"q2","options":["a","b"],"answer":0}]`)
	seed := []interface{}{
		&model.User{ID: "u1", Username: "u1", Nickname: "学习者", Role: "learner", Status: model.UserStatusActive, JoinedAt: now},
		&model.User{ID: "u2", Username: "u2", Nickname: "旁观者", Role: "learner", Status: model.UserStatusActive, JoinedAt: now},
		&model.User{ID: "co", Username: "co", Nickname: "企业", Role: "enterprise", Status: model.UserStatusActive, JoinedAt: now},
		&model.Career{ID: "c1", Name: "软件工程师"},
		&model.JobRole{ID: "r1", CareerID: "c1", Name: "后端工程师"},
		&model.RoleCapability{ID: "cap1", RoleID: "r1", Name: "并发能力", Weight: 1},
		&model.CapabilitySkill{ID: "cs1", CapabilityID: "cap1", SkillID: "s1", RequiredLevel: 3, Weight: 1},
		&model.Skill{ID: "s1", Name: "并发编程", Tags: datatypes.JSON(`[]`)},
		&model.KnowledgePoint{ID: "k0", Name: "线程基础", Difficulty: 1, EstimatedMinutes: 20, Questions: q, Tags: datatypes.JSON(`[]`)},
		&model.KnowledgePoint{ID: "k1", Name: "锁", Difficulty: 3, EstimatedMinutes: 30, Questions: q, Tags: datatypes.JSON(`[]`)},
		&model.SkillKnowledge{ID: "sk0", SkillID: "s1", KnowledgeID: "k0", Weight: 1},
		&model.SkillKnowledge{ID: "sk1", SkillID: "s1", KnowledgeID: "k1", Weight: 1},
		&model.KnowledgeRelation{ID: "rel1", FromID: "k0", ToID: "k1", Type: model.RelationPrerequisite},
	}
	for _, v := range seed {
		if err := db.Create(v).Error; err != nil {
			t.Fatalf("写入种子失败: %v", err)
		}
	}
	users := repository.NewUserRepo(db)
	gr := repository.NewGrowthRepo(db)
	marketRepo := repository.NewMarketRepo(db)
	catalog := NewCatalogService(repository.NewCatalogRepo(db), repository.NewCourseRepo(db))
	growth := NewGrowthService(catalog, gr, repository.NewProjectRepo(db), marketRepo, users, repository.NewKernelRepo(db))
	return &growthFixture{db: db, growth: growth, market: NewMarketService(marketRepo, growth, catalog, users, gr), gr: gr}
}

func TestGrowthLoopFromGoalToIncome(t *testing.T) {
	f := newGrowthFixture(t)

	if _, err := f.growth.SetGoal("u1", model.GoalRequest{RoleID: "r1", WeeklyHours: 5}); err != nil {
		t.Fatalf("设置目标失败: %v", err)
	}
	gap, err := f.growth.Gap("u1", "")
	if err != nil || gap.Readiness != 0 || len(gap.Gaps) != 1 || gap.Gaps[0].SkillID != "s1" {
		t.Fatalf("初始差距不正确: %+v %v", gap, err)
	}

	recs, _ := f.growth.Recommend("u1", gap, 0)
	if len(recs) != 2 || recs[0].Knowledge.ID != "k0" || !recs[0].Ready || recs[1].Ready {
		t.Fatalf("应先推荐前置知识 k0，k1 被阻塞: %+v", recs)
	}

	for i := 0; i < 6; i++ {
		for _, k := range []string{"k0", "k1"} {
			if _, err := f.growth.Practice("u1", k, []int{1, 0}); err != nil {
				t.Fatalf("练习失败: %v", err)
			}
		}
	}
	if lv := f.growth.SkillLevels("u1")["s1"]; lv != 2 {
		t.Fatalf("只有知识时技能应为 L2, got L%d", lv)
	}

	opp, err := f.market.Publish("co", model.OpportunityRequest{
		Title: "并发模块", Description: "开发并发模块", Budget: 3000,
		Requirements: []model.SkillRequirement{{SkillID: "s1", Level: 2}},
	})
	if err != nil {
		t.Fatalf("发布任务失败: %v", err)
	}
	app, err := f.market.Apply("u1", "learner", opp.ID, model.ApplyRequest{})
	if err != nil || app.MatchScore != 100 {
		t.Fatalf("申请失败或匹配度错误: %+v %v", app, err)
	}
	if _, err := f.market.Apply("u1", "learner", opp.ID, model.ApplyRequest{}); err == nil {
		t.Fatal("重复申请应被拒绝")
	}
	if _, err := f.market.Decide("u2", app.ID, false, model.ApplicationDecision{Action: "accept"}); err == nil {
		t.Fatal("非发布者不能处理申请")
	} else {
		assertStatus(t, err, http.StatusForbidden)
	}
	if _, err := f.market.Decide("co", app.ID, false, model.ApplicationDecision{Action: "complete"}); err == nil {
		t.Fatal("未录用的申请不能直接验收")
	}
	if _, err := f.market.Decide("co", app.ID, false, model.ApplicationDecision{Action: "accept"}); err != nil {
		t.Fatalf("录用失败: %v", err)
	}
	if _, err := f.market.Decide("co", app.ID, false, model.ApplicationDecision{Action: "complete", Rating: 5}); err != nil {
		t.Fatalf("验收失败: %v", err)
	}

	income, _ := f.gr.ListIncome("u1")
	if len(income) != 1 || income[0].Amount != 3000 {
		t.Fatalf("应生成一笔收入记录: %+v", income)
	}
	if lv := f.growth.SkillLevels("u1")["s1"]; lv < 3 {
		t.Fatalf("真实任务验收后技能应至少 L3, got L%d", lv)
	}
	gap, _ = f.growth.Gap("u1", "")
	if gap.Readiness != 100 || len(gap.Gaps) != 0 {
		t.Fatalf("岗位要求应全部达成: %+v", gap)
	}
	flywheel := f.growth.Flywheel("u1", 0)
	if flywheel.TasksCompleted != 1 || flywheel.Income != 3000 || flywheel.EvidenceCount != 1 {
		t.Fatalf("成长飞轮数据不正确: %+v", flywheel)
	}
}

func TestAddRelationRejectsCycle(t *testing.T) {
	f := newGrowthFixture(t)
	catalog := f.growth.catalog
	if _, err := catalog.AddRelation(model.KnowledgeRelation{FromID: "k1", ToID: "k0", Type: model.RelationPrerequisite}); err == nil {
		t.Fatal("k0→k1 已是前置关系，反向前置应形成环路被拒绝")
	}
	if _, err := catalog.AddRelation(model.KnowledgeRelation{FromID: "k1", ToID: "k0", Type: model.RelationRelated}); err != nil {
		t.Fatalf("关联关系不应受环路限制: %v", err)
	}
}
