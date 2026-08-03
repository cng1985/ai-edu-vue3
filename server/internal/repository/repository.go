package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }
type CourseRepo struct{ db *gorm.DB }
type QuizRepo struct{ db *gorm.DB }
type ReviewRepo struct{ db *gorm.DB }
type RoleRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo       { return &UserRepo{db: db} }
func NewCourseRepo(db *gorm.DB) *CourseRepo   { return &CourseRepo{db: db} }
func NewQuizRepo(db *gorm.DB) *QuizRepo       { return &QuizRepo{db: db} }
func NewReviewRepo(db *gorm.DB) *ReviewRepo   { return &ReviewRepo{db: db} }
func NewRoleRepo(db *gorm.DB) *RoleRepo       { return &RoleRepo{db: db} }

var Module = fx.Provide(NewUserRepo, NewCourseRepo, NewQuizRepo, NewReviewRepo, NewRoleRepo, NewSettingsRepo, NewCustomerRepo, NewDocumentRepo, NewKnowledgeRepo, NewAiModelRepo)

// --- User ---

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("LOWER(username) = LOWER(?)", username).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepo) List(keyword, role, status string, page, pageSize int) ([]model.User, int64, error) {
	q := r.db.Model(&model.User{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(username) LIKE LOWER(?) OR LOWER(nickname) LIKE LOWER(?)", kw, kw)
	}
	if role != "" {
		q = q.Where("role = ?", role)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []model.User
	offset := (page - 1) * pageSize
	err := q.Order("joined_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepo) Create(user *model.User) error { return r.db.Create(user).Error }
func (r *UserRepo) Update(user *model.User) error { return r.db.Save(user).Error }
func (r *UserRepo) Delete(id string) error        { return r.db.Delete(&model.User{}, "id = ?", id).Error }

func (r *UserRepo) CountByRole(role string) (int64, error) {
	var n int64
	err := r.db.Model(&model.User{}).Where("role = ?", role).Count(&n).Error
	return n, err
}

func (r *UserRepo) CountAdmins() (int64, error) {
	var n int64
	err := r.db.Model(&model.User{}).Where("role IN ?", []string{"admin", "reviewer", "operator"}).Count(&n).Error
	return n, err
}

func (r *UserRepo) CountActive() (int64, error) {
	var n int64
	err := r.db.Model(&model.User{}).Where("status = ?", "active").Count(&n).Error
	return n, err
}

func (r *UserRepo) Total() (int64, error) {
	var n int64
	err := r.db.Model(&model.User{}).Count(&n).Error
	return n, err
}

// --- Course ---

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

// --- Quiz ---

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

// --- Review ---

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

func (r *ReviewRepo) Update(review *model.Review) error { return r.db.Save(review).Error }

func (r *ReviewRepo) CountByStatus(status string) (int64, error) {
	var n int64
	err := r.db.Model(&model.Review{}).Where("status = ?", status).Count(&n).Error
	return n, err
}

func (r *ReviewRepo) Create(review *model.Review) error { return r.db.Create(review).Error }

// --- Role ---

func (r *RoleRepo) List() ([]model.RolePermission, error) {
	var roles []model.RolePermission
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByRole(role string) (*model.RolePermission, error) {
	var rp model.RolePermission
	err := r.db.First(&rp, "role = ?", role).Error
	return &rp, err
}

func (r *RoleRepo) Upsert(rp *model.RolePermission) error {
	return r.db.Save(rp).Error
}
