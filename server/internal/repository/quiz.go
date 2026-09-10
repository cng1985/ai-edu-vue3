package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type QuizRepo struct{ db *gorm.DB }

func NewQuizRepo(db *gorm.DB) *QuizRepo { return &QuizRepo{db: db} }

func (r *QuizRepo) List(keyword, courseID, status string, page, pageSize int) ([]model.Quiz, int64, error) {
	q := r.db.Model(&model.Quiz{})
	if keyword != "" {
		q = q.Where("LOWER(title) LIKE LOWER(?)", "%"+keyword+"%")
	}
	if courseID != "" {
		q = q.Where("course_id = ?", courseID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var quizzes []model.Quiz
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&quizzes).Error
	return quizzes, total, err
}

func (r *QuizRepo) FindByID(id string) (*model.Quiz, error) {
	var quiz model.Quiz
	err := r.db.First(&quiz, "id = ?", id).Error
	return &quiz, err
}

func (r *QuizRepo) Exists(id string) bool {
	var n int64
	r.db.Model(&model.Quiz{}).Where("id = ?", id).Count(&n)
	return n > 0
}

func (r *QuizRepo) Create(quiz *model.Quiz) error { return r.db.Create(quiz).Error }
func (r *QuizRepo) Update(quiz *model.Quiz) error { return r.db.Save(quiz).Error }
func (r *QuizRepo) Delete(id string) error        { return r.db.Delete(&model.Quiz{}, "id = ?", id).Error }

func (r *QuizRepo) Total() (int64, error) {
	var n int64
	err := r.db.Model(&model.Quiz{}).Count(&n).Error
	return n, err
}
