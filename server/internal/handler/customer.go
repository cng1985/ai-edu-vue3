package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/internal/ws"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CustomerHandler struct {
	svc *service.CustomerService
	hub *ws.Hub
	jwt *authutil.JWTManager
}

func NewCustomerHandler(svc *service.CustomerService, hub *ws.Hub, jwt *authutil.JWTManager) *CustomerHandler {
	return &CustomerHandler{svc: svc, hub: hub, jwt: jwt}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// --- REST: 学员端 ---

func (h *CustomerHandler) CreateTicket(c *gin.Context) {
	var req struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入咨询内容")
		return
	}
	claims := middleware.GetClaims(c)
	ticket, msg, err := h.svc.CreateTicket(claims.ID, req.Subject, req.Content)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, gin.H{"ticket": ticket, "message": msg}, "工单创建成功")
}

func (h *CustomerHandler) ListMyTickets(c *gin.Context) {
	claims := middleware.GetClaims(c)
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListMyTickets(claims.ID, page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CustomerHandler) GetTicket(c *gin.Context) {
	claims := middleware.GetClaims(c)
	ticket, err := h.svc.GetTicket(c.Param("id"), claims.ID, claims.Role)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, ticket)
}

func (h *CustomerHandler) ListMessages(c *gin.Context) {
	claims := middleware.GetClaims(c)
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListMessages(c.Param("id"), claims.ID, claims.Role, page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CustomerHandler) SendMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		response.Fail(c, http.StatusBadRequest, 400, "消息内容不能为空")
		return
	}
	claims := middleware.GetClaims(c)
	msg, err := h.svc.SendMessage(c.Param("id"), claims.ID, claims.Role, req.Content)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, msg, "发送成功")
}

// --- REST: 管理端 ---

func (h *CustomerHandler) AdminListTickets(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListTickets(c.Query("keyword"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CustomerHandler) AdminGetTicket(c *gin.Context) {
	claims := middleware.GetClaims(c)
	ticket, err := h.svc.GetTicket(c.Param("id"), claims.ID, claims.Role)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, ticket)
}

func (h *CustomerHandler) AdminListMessages(c *gin.Context) {
	claims := middleware.GetClaims(c)
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListMessages(c.Param("id"), claims.ID, claims.Role, page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CustomerHandler) AdminReply(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		response.Fail(c, http.StatusBadRequest, 400, "回复内容不能为空")
		return
	}
	claims := middleware.GetClaims(c)
	msg, err := h.svc.SendMessage(c.Param("id"), claims.ID, claims.Role, req.Content)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, msg, "回复成功")
}

func (h *CustomerHandler) AdminUpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请指定状态")
		return
	}
	claims := middleware.GetClaims(c)
	ticket, err := h.svc.UpdateStatus(c.Param("id"), req.Status, claims.ID, claims.Role)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, ticket, "状态已更新")
}

func (h *CustomerHandler) AdminStats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, stats)
}

// --- WebSocket ---

func (h *CustomerHandler) HandleWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		header := c.GetHeader("Authorization")
		if len(header) > 7 && header[:7] == "Bearer " {
			token = header[7:]
		}
	}
	if token == "" {
		response.Fail(c, http.StatusUnauthorized, 401, "未登录")
		return
	}
	claims, err := h.jwt.Verify(token)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, "登录已过期")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		UserID:   claims.ID,
		Username: claims.Username,
		Role:     claims.Role,
		IsAdmin:  rbac.IsAdminRole(claims.Role),
		Conn:     conn,
		Send:     make(chan []byte, 64),
		Tickets:  make(map[string]bool),
	}
	h.hub.Register(client)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *CustomerHandler) writePump(client *ws.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *CustomerHandler) readPump(client *ws.Client) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(8192)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		var env ws.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			continue
		}
		h.handleWSEvent(client, env)
	}
}

func (h *CustomerHandler) handleWSEvent(client *ws.Client, env ws.Envelope) {
	switch env.Event {
	case "support.subscribe":
		var req struct {
			TicketID string `json:"ticketId"`
		}
		if json.Unmarshal(env.Payload, &req) != nil || req.TicketID == "" {
			return
		}
		if _, err := h.svc.GetTicket(req.TicketID, client.UserID, client.Role); err != nil {
			return
		}
		h.hub.Subscribe(client, req.TicketID)
		if data, err := ws.MarshalEvent("support.subscribed", gin.H{"ticketId": req.TicketID}); err == nil {
			select {
			case client.Send <- data:
			default:
			}
		}

	case "support.send":
		var req struct {
			TicketID string `json:"ticketId"`
			Content  string `json:"content"`
		}
		if json.Unmarshal(env.Payload, &req) != nil || req.TicketID == "" || req.Content == "" {
			return
		}
		_, _ = h.svc.SendMessage(req.TicketID, client.UserID, client.Role, req.Content)
	}
}
