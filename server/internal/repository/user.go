package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

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

// CountAdmins 统计可登录管理后台的账号数
func (r *UserRepo) CountAdmins() (int64, error) {
	var n int64
	err := r.db.Model(&model.User{}).Where("role IN ?", rbac.AdminPortalRoles).Count(&n).Error
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
