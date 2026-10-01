package database

import (
	"os"
	"path/filepath"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/glebarez/sqlite"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return db, db.AutoMigrate(
		&model.User{},
		&model.Course{},
		&model.Chapter{},
		&model.Quiz{},
		&model.Review{},
		&model.RolePermission{},
		&model.SystemSetting{},
		&model.CustomerTicket{},
		&model.CustomerMessage{},
		&model.KnowledgeChunk{},
		&model.KnowledgeIndexMeta{},
		&model.CanonicalModel{},
		&model.Capability{},
		&model.CapabilityModel{},
		&model.Provider{},
		&model.ProviderModel{},
		&model.VirtualModel{},
		&model.VirtualModelMapping{},
		&model.Career{},
		&model.JobRole{},
		&model.RoleCapability{},
		&model.CapabilitySkill{},
		&model.Skill{},
		&model.SkillKnowledge{},
		&model.KnowledgePoint{},
		&model.KnowledgeRelation{},
		&model.ChapterKnowledge{},
		&model.StudentKnowledgeState{},
		&model.LearningEvent{},
		&model.StudentSkillState{},
		&model.UserGoal{},
		&model.Project{},
		&model.ProjectTask{},
		&model.TaskSubmission{},
		&model.Evidence{},
		&model.Opportunity{},
		&model.OpportunityApplication{},
		&model.IncomeRecord{},
		&model.CommunityPost{},
		&model.CommunityAnswer{},
		&model.KnowledgeResource{},
		&model.PipelineRun{},
	)
}

var Module = fx.Provide(NewDB)
