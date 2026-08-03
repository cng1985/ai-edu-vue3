package service

import (
	"errors"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/internal/ws"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

type CustomerService struct {
	customers *repository.CustomerRepo
	users     *repository.UserRepo
	hub       *ws.Hub
}

func NewCustomerService(customers *repository.CustomerRepo, users *repository.UserRepo, hub *ws.Hub) *CustomerService {
	return &CustomerService{customers: customers, users: users, hub: hub}
}

func (s *CustomerService) enrichTicket(ticket *model.CustomerTicket) {
	if user, err := s.users.FindByID(ticket.UserID); err == nil {
		ticket.UserNickname = user.Nickname
		ticket.UserUsername = user.Username
	}
	if msg, err := s.customers.GetLastMessage(ticket.ID); err == nil {
		ticket.LastMessage = msg.Content
	}
}

func (s *CustomerService) enrichMessage(msg *model.CustomerMessage) {
	if user, err := s.users.FindByID(msg.SenderID); err == nil {
		msg.SenderNickname = user.Nickname
	}
}

func (s *CustomerService) CreateTicket(userID, subject, content string) (*model.CustomerTicket, *model.CustomerMessage, error) {
	if subject == "" {
		subject = "咨询求助"
	}
	now := time.Now().UnixMilli()
	ticket := &model.CustomerTicket{
		ID: genID("ticket"), UserID: userID, Subject: subject,
		Status: "open", LastMessageAt: now, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.customers.CreateTicket(ticket); err != nil {
		return nil, nil, err
	}
	msg, err := s.sendMessage(ticket, userID, "learner", content)
	if err != nil {
		return nil, nil, err
	}
	s.enrichTicket(ticket)
	s.broadcastTicketNew(ticket)
	return ticket, msg, nil
}

func (s *CustomerService) ListMyTickets(userID string, page, pageSize int) (*model.PageResult[model.CustomerTicket], error) {
	tickets, total, err := s.customers.ListTicketsByUser(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	for i := range tickets {
		s.enrichTicket(&tickets[i])
	}
	return &model.PageResult[model.CustomerTicket]{List: tickets, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *CustomerService) ListTickets(keyword, status string, page, pageSize int) (*model.PageResult[model.CustomerTicket], error) {
	tickets, total, err := s.customers.ListTickets(keyword, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	for i := range tickets {
		s.enrichTicket(&tickets[i])
	}
	return &model.PageResult[model.CustomerTicket]{List: tickets, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *CustomerService) GetTicket(ticketID, userID, role string) (*model.CustomerTicket, error) {
	ticket, err := s.customers.FindTicketByID(ticketID)
	if err != nil {
		return nil, errors.New("工单不存在")
	}
	if !rbac.IsAdminRole(role) && ticket.UserID != userID {
		return nil, errors.New("无权限")
	}
	s.enrichTicket(ticket)
	return ticket, nil
}

func (s *CustomerService) ListMessages(ticketID, userID, role string, page, pageSize int) (*model.PageResult[model.CustomerMessage], error) {
	if _, err := s.GetTicket(ticketID, userID, role); err != nil {
		return nil, err
	}
	messages, total, err := s.customers.ListMessages(ticketID, page, pageSize)
	if err != nil {
		return nil, err
	}
	for i := range messages {
		s.enrichMessage(&messages[i])
	}
	return &model.PageResult[model.CustomerMessage]{List: messages, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *CustomerService) SendMessage(ticketID, senderID, senderRole, content string) (*model.CustomerMessage, error) {
	if content == "" {
		return nil, errors.New("消息内容不能为空")
	}
	ticket, err := s.customers.FindTicketByID(ticketID)
	if err != nil {
		return nil, errors.New("工单不存在")
	}
	if !rbac.IsAdminRole(senderRole) && ticket.UserID != senderID {
		return nil, errors.New("无权限")
	}
	if ticket.Status == "closed" {
		return nil, errors.New("工单已关闭，无法发送消息")
	}
	msg, err := s.sendMessage(ticket, senderID, senderRole, content)
	if err != nil {
		return nil, err
	}
	if rbac.IsAdminRole(senderRole) && ticket.Status == "open" {
		ticket.Status = "pending"
		ticket.UpdatedAt = time.Now().UnixMilli()
		_ = s.customers.UpdateTicket(ticket)
		s.broadcastTicketUpdate(ticket)
	}
	return msg, nil
}

func (s *CustomerService) sendMessage(ticket *model.CustomerTicket, senderID, senderRole, content string) (*model.CustomerMessage, error) {
	now := time.Now().UnixMilli()
	msg := &model.CustomerMessage{
		ID: genID("msg"), TicketID: ticket.ID, SenderID: senderID,
		SenderRole: senderRole, Content: content, CreatedAt: now,
	}
	if err := s.customers.CreateMessage(msg); err != nil {
		return nil, err
	}
	ticket.LastMessageAt = now
	ticket.UpdatedAt = now
	_ = s.customers.UpdateTicket(ticket)
	s.enrichMessage(msg)
	s.broadcastMessage(msg)
	return msg, nil
}

func (s *CustomerService) UpdateStatus(ticketID, status, operatorID, role string) (*model.CustomerTicket, error) {
	if !rbac.IsAdminRole(role) {
		return nil, errors.New("无权限")
	}
	if status != "open" && status != "pending" && status != "closed" {
		return nil, errors.New("无效的状态")
	}
	ticket, err := s.customers.FindTicketByID(ticketID)
	if err != nil {
		return nil, errors.New("工单不存在")
	}
	ticket.Status = status
	ticket.UpdatedAt = time.Now().UnixMilli()
	if err := s.customers.UpdateTicket(ticket); err != nil {
		return nil, err
	}
	s.enrichTicket(ticket)
	s.broadcastTicketUpdate(ticket)
	return ticket, nil
}

func (s *CustomerService) Stats() (*model.CustomerTicketStats, error) {
	total, _ := s.customers.TotalTickets()
	open, _ := s.customers.CountByStatus("open")
	pending, _ := s.customers.CountByStatus("pending")
	closed, _ := s.customers.CountByStatus("closed")
	return &model.CustomerTicketStats{Total: total, Open: open, Pending: pending, Closed: closed}, nil
}

func (s *CustomerService) broadcastMessage(msg *model.CustomerMessage) {
	data, err := ws.MarshalEvent("support.message", msg)
	if err != nil {
		return
	}
	s.hub.BroadcastToTicket(msg.TicketID, data, nil)
}

func (s *CustomerService) broadcastTicketUpdate(ticket *model.CustomerTicket) {
	data, err := ws.MarshalEvent("support.ticket.update", ticket)
	if err != nil {
		return
	}
	s.hub.BroadcastToTicket(ticket.ID, data, nil)
	s.hub.SendToUser(ticket.UserID, data)
}

func (s *CustomerService) broadcastTicketNew(ticket *model.CustomerTicket) {
	data, err := ws.MarshalEvent("support.ticket.new", ticket)
	if err != nil {
		return
	}
	s.hub.BroadcastToAdmins(data)
}
