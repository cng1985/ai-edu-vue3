package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

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
