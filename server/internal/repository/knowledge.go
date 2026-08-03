package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type KnowledgeRepo struct{ db *gorm.DB }

func NewKnowledgeRepo(db *gorm.DB) *KnowledgeRepo { return &KnowledgeRepo{db: db} }

func (r *KnowledgeRepo) CountChunks() (int64, error) {
	var n int64
	err := r.db.Model(&model.KnowledgeChunk{}).Count(&n).Error
	return n, err
}

func (r *KnowledgeRepo) CountCourses() (int64, error) {
	var n int64
	err := r.db.Model(&model.KnowledgeChunk{}).Distinct("course_id").Count(&n).Error
	return n, err
}

func (r *KnowledgeRepo) CountChapters() (int64, error) {
	var n int64
	err := r.db.Model(&model.KnowledgeChunk{}).Distinct("chapter_id").Count(&n).Error
	return n, err
}

func (r *KnowledgeRepo) ListAll() ([]model.KnowledgeChunk, error) {
	var chunks []model.KnowledgeChunk
	err := r.db.Find(&chunks).Error
	return chunks, err
}

func (r *KnowledgeRepo) DeleteByChapter(courseID, chapterID string) error {
	return r.db.Delete(&model.KnowledgeChunk{}, "course_id = ? AND chapter_id = ?", courseID, chapterID).Error
}

func (r *KnowledgeRepo) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&model.KnowledgeChunk{}).Error
}

func (r *KnowledgeRepo) UpsertChunks(chunks []model.KnowledgeChunk) error {
	if len(chunks) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		courseID := chunks[0].CourseID
		chapterID := chunks[0].ChapterID
		if err := tx.Delete(&model.KnowledgeChunk{}, "course_id = ? AND chapter_id = ?", courseID, chapterID).Error; err != nil {
			return err
		}
		return tx.CreateInBatches(chunks, 50).Error
	})
}

func (r *KnowledgeRepo) GetMeta(key string) (string, error) {
	var meta model.KnowledgeIndexMeta
	err := r.db.First(&meta, "key = ?", key).Error
	if err != nil {
		return "", err
	}
	return meta.Value, nil
}

func (r *KnowledgeRepo) SetMeta(key, value string, updatedAt int64) error {
	meta := model.KnowledgeIndexMeta{Key: key, Value: value, UpdatedAt: updatedAt}
	return r.db.Save(&meta).Error
}

func (r *KnowledgeRepo) List(page, pageSize int) ([]model.KnowledgeChunk, int64, error) {
	q := r.db.Model(&model.KnowledgeChunk{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var chunks []model.KnowledgeChunk
	err := q.Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&chunks).Error
	return chunks, total, err
}
