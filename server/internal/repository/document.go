package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type DocumentRepo struct{ db *gorm.DB }

func NewDocumentRepo(db *gorm.DB) *DocumentRepo { return &DocumentRepo{db: db} }

func (r *DocumentRepo) List(keyword, docType, status string, page, pageSize int) ([]model.Document, int64, error) {
	q := r.db.Model(&model.Document{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(doc_no) LIKE LOWER(?) OR LOWER(title) LIKE LOWER(?)", kw, kw)
	}
	if docType != "" {
		q = q.Where("type = ?", docType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var docs []model.Document
	offset := (page - 1) * pageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&docs).Error
	return docs, total, err
}

func (r *DocumentRepo) FindByID(id string) (*model.Document, error) {
	var doc model.Document
	err := r.db.First(&doc, "id = ?", id).Error
	return &doc, err
}

func (r *DocumentRepo) FindByDocNo(docNo string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Where("doc_no = ?", docNo).First(&doc).Error
	return &doc, err
}

func (r *DocumentRepo) Create(doc *model.Document) error { return r.db.Create(doc).Error }

func (r *DocumentRepo) Update(doc *model.Document) error { return r.db.Save(doc).Error }

func (r *DocumentRepo) Delete(id string) error { return r.db.Delete(&model.Document{}, "id = ?", id).Error }

func (r *DocumentRepo) ListAll(keyword, docType, status string) ([]model.Document, error) {
	q := r.db.Model(&model.Document{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("LOWER(doc_no) LIKE LOWER(?) OR LOWER(title) LIKE LOWER(?)", kw, kw)
	}
	if docType != "" {
		q = q.Where("type = ?", docType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var docs []model.Document
	err := q.Order("created_at DESC").Find(&docs).Error
	return docs, err
}
