package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type ReviewRepo struct{ db *gorm.DB }

func NewReviewRepo(db *gorm.DB) *ReviewRepo { return &ReviewRepo{db: db} }

func (r *ReviewRepo) List(status, reviewType string, page, pageSize int) ([]model.Review, int64, error) {
	q := r.db.Model(&model.Review{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if reviewType != "" {
		q = q.Where("type = ?", reviewType)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var reviews []model.Review
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepo) FindByID(id string) (*model.Review, error) {
	var review model.Review
	err := r.db.First(&review, "id = ?", id).Error
	return &review, err
}

func (r *ReviewRepo) Create(review *model.Review) error { return r.db.Create(review).Error }
func (r *ReviewRepo) Update(review *model.Review) error { return r.db.Save(review).Error }

func (r *ReviewRepo) CountByStatus(status string) (int64, error) {
	var n int64
	err := r.db.Model(&model.Review{}).Where("status = ?", status).Count(&n).Error
	return n, err
}
