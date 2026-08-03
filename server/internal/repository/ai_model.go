package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type AiModelRepo struct{ db *gorm.DB }

func NewAiModelRepo(db *gorm.DB) *AiModelRepo { return &AiModelRepo{db: db} }

// Transaction 在同一事务中执行一组模型配置变更。
func (r *AiModelRepo) Transaction(fn func(*AiModelRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewAiModelRepo(tx))
	})
}

// --- CanonicalModel ---

func (r *AiModelRepo) ListCanonicalModels(keyword string, page, pageSize int) ([]model.CanonicalModel, int64, error) {
	q := r.db.Model(&model.CanonicalModel{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(code) LIKE LOWER(?) OR LOWER(name) LIKE LOWER(?)", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CanonicalModel
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *AiModelRepo) FindCanonicalModel(id string) (*model.CanonicalModel, error) {
	var m model.CanonicalModel
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) FindCanonicalModelByCode(code string) (*model.CanonicalModel, error) {
	var m model.CanonicalModel
	err := r.db.Where("code = ?", code).First(&m).Error
	return &m, err
}

func (r *AiModelRepo) CreateCanonicalModel(m *model.CanonicalModel) error {
	return r.db.Create(m).Error
}
func (r *AiModelRepo) UpdateCanonicalModel(m *model.CanonicalModel) error { return r.db.Save(m).Error }
func (r *AiModelRepo) DeleteCanonicalModel(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CapabilityModel{}, "canonical_model_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.ProviderModel{}, "canonical_model_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.VirtualModelMapping{}, "canonical_model_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.CanonicalModel{}, "id = ?", id).Error
	})
}

func (r *AiModelRepo) CountCanonicalModels() (int64, error) {
	var n int64
	err := r.db.Model(&model.CanonicalModel{}).Count(&n).Error
	return n, err
}

// --- Capability ---

func (r *AiModelRepo) ListCapabilities(keyword string, page, pageSize int) ([]model.Capability, int64, error) {
	q := r.db.Model(&model.Capability{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(code) LIKE LOWER(?) OR LOWER(name) LIKE LOWER(?)", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Capability
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *AiModelRepo) FindCapability(id string) (*model.Capability, error) {
	var m model.Capability
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) CreateCapability(m *model.Capability) error { return r.db.Create(m).Error }
func (r *AiModelRepo) UpdateCapability(m *model.Capability) error { return r.db.Save(m).Error }
func (r *AiModelRepo) DeleteCapability(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CapabilityModel{}, "capability_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Capability{}, "id = ?", id).Error
	})
}

func (r *AiModelRepo) CountCapabilities() (int64, error) {
	var n int64
	err := r.db.Model(&model.Capability{}).Count(&n).Error
	return n, err
}

// --- CapabilityModel ---

func (r *AiModelRepo) ListCapabilityModels(canonicalModelID string) ([]model.CapabilityModel, error) {
	q := r.db.Model(&model.CapabilityModel{})
	if canonicalModelID != "" {
		q = q.Where("canonical_model_id = ?", canonicalModelID)
	}
	var list []model.CapabilityModel
	err := q.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *AiModelRepo) FindCapabilityModel(id string) (*model.CapabilityModel, error) {
	var m model.CapabilityModel
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) CreateCapabilityModel(m *model.CapabilityModel) error {
	return r.db.Create(m).Error
}
func (r *AiModelRepo) DeleteCapabilityModel(id string) error {
	return r.db.Delete(&model.CapabilityModel{}, "id = ?", id).Error
}

func (r *AiModelRepo) ExistsCapabilityModel(canonicalModelID, capabilityID string) bool {
	var n int64
	r.db.Model(&model.CapabilityModel{}).
		Where("canonical_model_id = ? AND capability_id = ?", canonicalModelID, capabilityID).
		Count(&n)
	return n > 0
}

// --- Provider ---

func (r *AiModelRepo) ListProviders(keyword string, page, pageSize int) ([]model.Provider, int64, error) {
	q := r.db.Model(&model.Provider{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(code) LIKE LOWER(?) OR LOWER(name) LIKE LOWER(?)", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Provider
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	for i := range list {
		var cnt int64
		r.db.Model(&model.ProviderModel{}).Where("provider_id = ?", list[i].ID).Count(&cnt)
		list[i].ModelCount = int(cnt)
	}
	return list, total, err
}

func (r *AiModelRepo) FindProvider(id string) (*model.Provider, error) {
	var m model.Provider
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) FindProviderByCode(code string) (*model.Provider, error) {
	var m model.Provider
	err := r.db.Where("code = ?", code).First(&m).Error
	return &m, err
}

func (r *AiModelRepo) CreateProvider(m *model.Provider) error { return r.db.Create(m).Error }
func (r *AiModelRepo) UpdateProvider(m *model.Provider) error { return r.db.Save(m).Error }
func (r *AiModelRepo) DeleteProvider(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.ProviderModel{}, "provider_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Provider{}, "id = ?", id).Error
	})
}

func (r *AiModelRepo) CountProviders() (int64, error) {
	var n int64
	err := r.db.Model(&model.Provider{}).Count(&n).Error
	return n, err
}

// --- ProviderModel ---

func (r *AiModelRepo) ListProviderModels(providerID, canonicalModelID string) ([]model.ProviderModel, error) {
	q := r.db.Model(&model.ProviderModel{})
	if providerID != "" {
		q = q.Where("provider_id = ?", providerID)
	}
	if canonicalModelID != "" {
		q = q.Where("canonical_model_id = ?", canonicalModelID)
	}
	var list []model.ProviderModel
	err := q.Order("priority ASC, created_at DESC").Find(&list).Error
	return list, err
}

func (r *AiModelRepo) FindProviderModel(id string) (*model.ProviderModel, error) {
	var m model.ProviderModel
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) CreateProviderModel(m *model.ProviderModel) error { return r.db.Create(m).Error }
func (r *AiModelRepo) UpdateProviderModel(m *model.ProviderModel) error { return r.db.Save(m).Error }
func (r *AiModelRepo) DeleteProviderModel(id string) error {
	return r.db.Delete(&model.ProviderModel{}, "id = ?", id).Error
}

func (r *AiModelRepo) CountProviderModels() (int64, error) {
	var n int64
	err := r.db.Model(&model.ProviderModel{}).Count(&n).Error
	return n, err
}

func (r *AiModelRepo) ListActiveProviderModelsByCanonical(canonicalModelID string) ([]model.ProviderModel, error) {
	var list []model.ProviderModel
	err := r.db.Where("canonical_model_id = ? AND status = 1", canonicalModelID).
		Order("priority ASC").Find(&list).Error
	return list, err
}

// --- VirtualModel ---

func (r *AiModelRepo) ListVirtualModels(keyword string, page, pageSize int) ([]model.VirtualModel, int64, error) {
	q := r.db.Model(&model.VirtualModel{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(code) LIKE LOWER(?) OR LOWER(name) LIKE LOWER(?)", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.VirtualModel
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	for i := range list {
		var cnt int64
		r.db.Model(&model.VirtualModelMapping{}).Where("virtual_model_id = ?", list[i].ID).Count(&cnt)
		list[i].MappingCount = int(cnt)
	}
	return list, total, err
}

func (r *AiModelRepo) FindVirtualModel(id string) (*model.VirtualModel, error) {
	var m model.VirtualModel
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) FindVirtualModelByCode(code string) (*model.VirtualModel, error) {
	var m model.VirtualModel
	err := r.db.Where("code = ?", code).First(&m).Error
	return &m, err
}

func (r *AiModelRepo) CreateVirtualModel(m *model.VirtualModel) error { return r.db.Create(m).Error }
func (r *AiModelRepo) UpdateVirtualModel(m *model.VirtualModel) error { return r.db.Save(m).Error }
func (r *AiModelRepo) DeleteVirtualModel(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.VirtualModelMapping{}, "virtual_model_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.VirtualModel{}, "id = ?", id).Error
	})
}

func (r *AiModelRepo) CountVirtualModels() (int64, error) {
	var n int64
	err := r.db.Model(&model.VirtualModel{}).Count(&n).Error
	return n, err
}

// --- VirtualModelMapping ---

func (r *AiModelRepo) ListVirtualModelMappings(virtualModelID string) ([]model.VirtualModelMapping, error) {
	q := r.db.Model(&model.VirtualModelMapping{})
	if virtualModelID != "" {
		q = q.Where("virtual_model_id = ?", virtualModelID)
	}
	var list []model.VirtualModelMapping
	err := q.Order("priority ASC, created_at DESC").Find(&list).Error
	return list, err
}

func (r *AiModelRepo) FindVirtualModelMapping(id string) (*model.VirtualModelMapping, error) {
	var m model.VirtualModelMapping
	err := r.db.First(&m, "id = ?", id).Error
	return &m, err
}

func (r *AiModelRepo) CreateVirtualModelMapping(m *model.VirtualModelMapping) error {
	return r.db.Create(m).Error
}

func (r *AiModelRepo) UpdateVirtualModelMapping(m *model.VirtualModelMapping) error {
	return r.db.Save(m).Error
}

func (r *AiModelRepo) DeleteVirtualModelMapping(id string) error {
	return r.db.Delete(&model.VirtualModelMapping{}, "id = ?", id).Error
}

func (r *AiModelRepo) ExistsVirtualModelMapping(virtualModelID, canonicalModelID string) bool {
	var n int64
	r.db.Model(&model.VirtualModelMapping{}).
		Where("virtual_model_id = ? AND canonical_model_id = ?", virtualModelID, canonicalModelID).
		Count(&n)
	return n > 0
}

func (r *AiModelRepo) ListActiveMappingsByVirtual(virtualModelID string) ([]model.VirtualModelMapping, error) {
	var list []model.VirtualModelMapping
	err := r.db.Where("virtual_model_id = ?", virtualModelID).
		Order("priority ASC").Find(&list).Error
	return list, err
}
