package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

// ProjectRepo 项目实践：项目、任务与提交。
type ProjectRepo struct{ db *gorm.DB }

func NewProjectRepo(db *gorm.DB) *ProjectRepo { return &ProjectRepo{db: db} }

func (r *ProjectRepo) List(status string) ([]model.Project, error) {
	var list []model.Project
	q := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).Order("created_at ASC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *ProjectRepo) Find(id string) (*model.Project, error) {
	var p model.Project
	err := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).First(&p, "id = ?", id).Error
	return &p, err
}

func (r *ProjectRepo) FindTask(projectID, taskID string) (*model.ProjectTask, error) {
	var t model.ProjectTask
	err := r.db.Where("project_id = ? AND id = ?", projectID, taskID).First(&t).Error
	return &t, err
}

func (r *ProjectRepo) Save(v interface{}) error { return r.db.Save(v).Error }

func (r *ProjectRepo) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.ProjectTask{}, "project_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Project{}, "id = ?", id).Error
	})
}

func (r *ProjectRepo) DeleteTask(projectID, taskID string) error {
	return r.db.Delete(&model.ProjectTask{}, "project_id = ? AND id = ?", projectID, taskID).Error
}

func (r *ProjectRepo) AddSubmission(s *model.TaskSubmission) error { return r.db.Create(s).Error }

func (r *ProjectRepo) SaveSubmission(s *model.TaskSubmission) error { return r.db.Save(s).Error }

func (r *ProjectRepo) ListSubmissions(userID, projectID string, limit int) ([]model.TaskSubmission, error) {
	var list []model.TaskSubmission
	q := r.db.Order("created_at DESC")
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if projectID != "" {
		q = q.Where("project_id = ?", projectID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *ProjectRepo) Count() int64 {
	var n int64
	r.db.Model(&model.Project{}).Count(&n)
	return n
}

func (r *ProjectRepo) CountSubmissions() int64 {
	var n int64
	r.db.Model(&model.TaskSubmission{}).Count(&n)
	return n
}

// MarketRepo IT 任务市场：企业任务与申请。
type MarketRepo struct{ db *gorm.DB }

func NewMarketRepo(db *gorm.DB) *MarketRepo { return &MarketRepo{db: db} }

func (r *MarketRepo) List(status, publisherID string) ([]model.Opportunity, error) {
	var list []model.Opportunity
	q := r.db.Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if publisherID != "" {
		q = q.Where("publisher_id = ?", publisherID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *MarketRepo) Find(id string) (*model.Opportunity, error) {
	var o model.Opportunity
	err := r.db.First(&o, "id = ?", id).Error
	return &o, err
}

func (r *MarketRepo) Save(o *model.Opportunity) error { return r.db.Save(o).Error }

func (r *MarketRepo) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.OpportunityApplication{}, "opportunity_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Opportunity{}, "id = ?", id).Error
	})
}

func (r *MarketRepo) FindApplication(id string) (*model.OpportunityApplication, error) {
	var a model.OpportunityApplication
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *MarketRepo) FindUserApplication(opportunityID, userID string) (*model.OpportunityApplication, error) {
	var a model.OpportunityApplication
	err := r.db.Where("opportunity_id = ? AND user_id = ?", opportunityID, userID).First(&a).Error
	return &a, err
}

func (r *MarketRepo) SaveApplication(a *model.OpportunityApplication) error { return r.db.Save(a).Error }

func (r *MarketRepo) ListApplications(opportunityID, userID string) ([]model.OpportunityApplication, error) {
	var list []model.OpportunityApplication
	q := r.db.Order("created_at DESC")
	if opportunityID != "" {
		q = q.Where("opportunity_id = ?", opportunityID)
	}
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *MarketRepo) CountByStatus(status string) int64 {
	var n int64
	q := r.db.Model(&model.Opportunity{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&n)
	return n
}

func (r *MarketRepo) CountApplications(status string) int64 {
	var n int64
	q := r.db.Model(&model.OpportunityApplication{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&n)
	return n
}
