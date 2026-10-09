package handler

import (
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/pipeline"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// KernelHandler AI 学习内核与 Agent 接口。
type KernelHandler struct {
	kernel *service.KernelService
	agents *service.AgentService
}

func NewKernelHandler(kernel *service.KernelService, agents *service.AgentService) *KernelHandler {
	return &KernelHandler{kernel: kernel, agents: agents}
}

func (h *KernelHandler) Agents(c *gin.Context) {
	response.OK(c, h.agents.List())
}

func (h *KernelHandler) AgentContext(c *gin.Context) {
	claims := middleware.GetClaims(c)
	text, err := h.agents.ContextPreview(c.Request.Context(), claims.ID, claims.Role, c.Param("code"), c.Query("message"))
	reply(c, gin.H{"context": text}, err)
}

// AgentChatStream 与 Agent 进行 SSE 流式对话。
func (h *KernelHandler) AgentChatStream(c *gin.Context) {
	var req model.AgentChatRequest
	if !bindJSON(c, &req) {
		return
	}
	flusher, ok := sseWriter(c)
	if !ok {
		return
	}
	claims := middleware.GetClaims(c)
	res, err := h.agents.Chat(c.Request.Context(), claims.ID, claims.Role, c.Param("code"), req, func(token string) {
		sseSend(c, flusher, gin.H{"type": "token", "content": token})
	})
	if err != nil {
		sseSend(c, flusher, gin.H{"type": "error", "message": err.Error()})
		return
	}
	sseSend(c, flusher, gin.H{"type": "done", "text": res.Text, "sources": res.Sources})
}

func (h *KernelHandler) Stages(c *gin.Context) {
	response.OK(c, h.kernel.Stages())
}

// RunStream 运行学习 Pipeline，按 Stage 推送执行轨迹，AI 回答阶段推送 token。
func (h *KernelHandler) RunStream(c *gin.Context) {
	var req model.KernelRunRequest
	_ = c.ShouldBindJSON(&req)
	flusher, ok := sseWriter(c)
	if !ok {
		return
	}
	run, err := h.kernel.RunLearning(c.Request.Context(), middleware.GetClaims(c).ID, req.Question,
		func(tr pipeline.Trace) { sseSend(c, flusher, gin.H{"type": "stage", "stage": tr}) },
		func(token string) { sseSend(c, flusher, gin.H{"type": "token", "content": token}) },
	)
	if err != nil {
		sseSend(c, flusher, gin.H{"type": "error", "message": err.Error()})
		return
	}
	sseSend(c, flusher, gin.H{"type": "done", "run": run})
}

func (h *KernelHandler) Run(c *gin.Context) {
	var req model.KernelRunRequest
	_ = c.ShouldBindJSON(&req)
	run, err := h.kernel.RunLearning(c.Request.Context(), middleware.GetClaims(c).ID, req.Question, nil, nil)
	reply(c, run, err)
}

func (h *KernelHandler) MyRuns(c *gin.Context) {
	list, err := h.kernel.Runs(middleware.GetClaims(c).ID, intQuery(c, "limit", 20))
	reply(c, list, err)
}

func (h *KernelHandler) AllRuns(c *gin.Context) {
	list, err := h.kernel.Runs(c.Query("userId"), intQuery(c, "limit", 50))
	reply(c, list, err)
}
