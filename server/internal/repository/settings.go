package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type SettingsRepo struct{ db *gorm.DB }

func NewSettingsRepo(db *gorm.DB) *SettingsRepo { return &SettingsRepo{db: db} }

func (r *SettingsRepo) Get(key string) (string, error) {
	var row model.SystemSetting
	err := r.db.First(&row, "key = ?", key).Error
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func (r *SettingsRepo) GetAll() (map[string]string, error) {
	var rows []model.SystemSetting
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out, nil
}

func (r *SettingsRepo) List() ([]model.SystemSetting, error) {
	var rows []model.SystemSetting
	err := r.db.Find(&rows).Error
	return rows, err
}

func (r *SettingsRepo) Upsert(key, value string, updatedAt int64) error {
	return r.db.Save(&model.SystemSetting{
		Key: key, Value: value, UpdatedAt: updatedAt,
	}).Error
}

func (r *SettingsRepo) HasAny() (bool, error) {
	var n int64
	err := r.db.Model(&model.SystemSetting{}).Count(&n).Error
	return n > 0, err
}
