package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type CourseRepo struct{ db *gorm.DB }

func NewCourseRepo(db *gorm.DB) *CourseRepo { return &CourseRepo{db: db} }

func (r *CourseRepo) List(keyword, status string, page, pageSize int) ([]model.Course, int64, error) {
	q := r.db.Model(&model.Course{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(title) LIKE LOWER(?) OR id LIKE ?", kw, kw)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var courses []model.Course
	offset := (page - 1) * pageSize
	err := q.Preload("Chapters").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error
	for i := range courses {
		courses[i].ChapterCount = len(courses[i].Chapters)
	}
	return courses, total, err
}

func (r *CourseRepo) FindByID(id string) (*model.Course, error) {
	var course model.Course
	err := r.db.Preload("Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at ASC")
	}).First(&course, "id = ?", id).Error
	return &course, err
}

func (r *CourseRepo) Exists(id string) bool {
	var n int64
	r.db.Model(&model.Course{}).Where("id = ?", id).Count(&n)
	return n > 0
}

func (r *CourseRepo) Create(course *model.Course) error { return r.db.Create(course).Error }
func (r *CourseRepo) Update(course *model.Course) error { return r.db.Save(course).Error }

func (r *CourseRepo) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Chapter{}, "course_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Course{}, "id = ?", id).Error
	})
}

func (r *CourseRepo) AddChapter(ch *model.Chapter) error { return r.db.Create(ch).Error }

func (r *CourseRepo) FindChapter(courseID, chapterID string) (*model.Chapter, error) {
	var ch model.Chapter
	err := r.db.Where("course_id = ? AND id = ?", courseID, chapterID).First(&ch).Error
	return &ch, err
}

func (r *CourseRepo) UpdateChapter(ch *model.Chapter) error { return r.db.Save(ch).Error }

func (r *CourseRepo) DeleteChapter(courseID, chapterID string) error {
	return r.db.Delete(&model.Chapter{}, "course_id = ? AND id = ?", courseID, chapterID).Error
}

func (r *CourseRepo) ChapterExists(courseID, chapterID string) bool {
	var n int64
	r.db.Model(&model.Chapter{}).Where("course_id = ? AND id = ?", courseID, chapterID).Count(&n)
	return n > 0
}

func (r *CourseRepo) CountByStatus(status string) (int64, error) {
	var n int64
	err := r.db.Model(&model.Course{}).Where("status = ?", status).Count(&n).Error
	return n, err
}

func (r *CourseRepo) Total() (int64, error) {
	var n int64
	err := r.db.Model(&model.Course{}).Count(&n).Error
	return n, err
}

func (r *CourseRepo) TotalChapters() (int64, error) {
	var n int64
	err := r.db.Model(&model.Chapter{}).Count(&n).Error
	return n, err
}

func (r *CourseRepo) ListPublished() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Preload("Chapters", "status = ?", "published").
		Where("status = ?", "published").
		Find(&courses).Error
	return courses, err
}

func (r *CourseRepo) FindPublishedByID(id string) (*model.Course, error) {
	var course model.Course
	err := r.db.Preload("Chapters", "status = ?", "published").
		Where("id = ? AND status = ?", id, "published").
		First(&course).Error
	return &course, err
}
