package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

type CustomerRepo struct{ db *gorm.DB }

func NewCustomerRepo(db *gorm.DB) *CustomerRepo { return &CustomerRepo{db: db} }

func (r *CustomerRepo) CreateTicket(ticket *model.CustomerTicket) error {
	return r.db.Create(ticket).Error
}

func (r *CustomerRepo) UpdateTicket(ticket *model.CustomerTicket) error {
	return r.db.Save(ticket).Error
}

func (r *CustomerRepo) FindTicketByID(id string) (*model.CustomerTicket, error) {
	var ticket model.CustomerTicket
	err := r.db.First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *CustomerRepo) ListTicketsByUser(userID string, page, pageSize int) ([]model.CustomerTicket, int64, error) {
	q := r.db.Model(&model.CustomerTicket{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var tickets []model.CustomerTicket
	offset := (page - 1) * pageSize
	err := q.Order("last_message_at DESC").Offset(offset).Limit(pageSize).Find(&tickets).Error
	return tickets, total, err
}

func (r *CustomerRepo) ListTickets(keyword, status string, page, pageSize int) ([]model.CustomerTicket, int64, error) {
	q := r.db.Model(&model.CustomerTicket{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("subject LIKE ? OR user_id IN (SELECT id FROM users WHERE LOWER(username) LIKE LOWER(?) OR LOWER(nickname) LIKE LOWER(?))", kw, kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var tickets []model.CustomerTicket
	offset := (page - 1) * pageSize
	err := q.Order("last_message_at DESC").Offset(offset).Limit(pageSize).Find(&tickets).Error
	return tickets, total, err
}

func (r *CustomerRepo) CreateMessage(msg *model.CustomerMessage) error {
	return r.db.Create(msg).Error
}

func (r *CustomerRepo) ListMessages(ticketID string, page, pageSize int) ([]model.CustomerMessage, int64, error) {
	q := r.db.Model(&model.CustomerMessage{}).Where("ticket_id = ?", ticketID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var messages []model.CustomerMessage
	offset := (page - 1) * pageSize
	err := q.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&messages).Error
	return messages, total, err
}

func (r *CustomerRepo) GetLastMessage(ticketID string) (*model.CustomerMessage, error) {
	var msg model.CustomerMessage
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at DESC").First(&msg).Error
	return &msg, err
}

func (r *CustomerRepo) CountByStatus(status string) (int64, error) {
	var n int64
	err := r.db.Model(&model.CustomerTicket{}).Where("status = ?", status).Count(&n).Error
	return n, err
}

func (r *CustomerRepo) TotalTickets() (int64, error) {
	var n int64
	err := r.db.Model(&model.CustomerTicket{}).Count(&n).Error
	return n, err
}
