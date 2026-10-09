package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GrowthRepo 个人成长数据：目标、知识状态、学习事件、技能状态、证据与收入。
type GrowthRepo struct{ db *gorm.DB }

func NewGrowthRepo(db *gorm.DB) *GrowthRepo { return &GrowthRepo{db: db} }

func (r *GrowthRepo) FindGoal(userID string) (*model.UserGoal, error) {
	var g model.UserGoal
	err := r.db.First(&g, "user_id = ?", userID).Error
	return &g, err
}

func (r *GrowthRepo) SaveGoal(g *model.UserGoal) error { return r.db.Save(g).Error }

func (r *GrowthRepo) ListKnowledgeStates(userID string) ([]model.StudentKnowledgeState, error) {
	var list []model.StudentKnowledgeState
	err := r.db.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *GrowthRepo) FindKnowledgeState(userID, knowledgeID string) (*model.StudentKnowledgeState, error) {
	var s model.StudentKnowledgeState
	err := r.db.Where("user_id = ? AND knowledge_id = ?", userID, knowledgeID).First(&s).Error
	return &s, err
}

func (r *GrowthRepo) SaveKnowledgeState(s *model.StudentKnowledgeState) error {
	return r.db.Save(s).Error
}

func (r *GrowthRepo) AddEvent(e *model.LearningEvent) error { return r.db.Create(e).Error }

func (r *GrowthRepo) ListEvents(userID string, since int64) ([]model.LearningEvent, error) {
	var list []model.LearningEvent
	err := r.db.Where("user_id = ? AND created_at >= ?", userID, since).Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *GrowthRepo) CountEvents(userID string) int64 {
	var n int64
	r.db.Model(&model.LearningEvent{}).Where("user_id = ?", userID).Count(&n)
	return n
}

func (r *GrowthRepo) ListSkillStates(userID string) ([]model.StudentSkillState, error) {
	var list []model.StudentSkillState
	err := r.db.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// ListAllSkillStates 用于人才库检索与任务匹配。
func (r *GrowthRepo) ListAllSkillStates() ([]model.StudentSkillState, error) {
	var list []model.StudentSkillState
	err := r.db.Where("level > 0").Find(&list).Error
	return list, err
}

// UpsertSkillState 按 (user, skill) 唯一键写入技能状态。
func (r *GrowthRepo) UpsertSkillState(s *model.StudentSkillState) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "skill_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"level", "score", "knowledge_score", "project_score", "task_score", "evidence_count", "updated_at"}),
	}).Create(s).Error
}

func (r *GrowthRepo) AddEvidence(e *model.Evidence) error { return r.db.Create(e).Error }

func (r *GrowthRepo) ListEvidence(userID string) ([]model.Evidence, error) {
	var list []model.Evidence
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *GrowthRepo) AddIncome(i *model.IncomeRecord) error { return r.db.Create(i).Error }

func (r *GrowthRepo) ListIncome(userID string) ([]model.IncomeRecord, error) {
	var list []model.IncomeRecord
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *GrowthRepo) SumIncome() float64 {
	var total float64
	r.db.Model(&model.IncomeRecord{}).Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total
}

func (r *GrowthRepo) CountGoals() int64 {
	var n int64
	r.db.Model(&model.UserGoal{}).Count(&n)
	return n
}

func (r *GrowthRepo) CountEvidence() int64 {
	var n int64
	r.db.Model(&model.Evidence{}).Count(&n)
	return n
}

func (r *GrowthRepo) CountAllEvents() int64 {
	var n int64
	r.db.Model(&model.LearningEvent{}).Count(&n)
	return n
}
