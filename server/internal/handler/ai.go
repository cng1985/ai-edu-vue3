package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AIHandler struct{ svc *service.AIService }

func NewAIHandler(svc *service.AIService) *AIHandler { return &AIHandler{svc: svc} }

func (h *AIHandler) Config(c *gin.Context) {
	response.OK(c, h.svc.ConfigInfo())
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Question == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入问题")
		return
	}
	result, err := h.svc.Chat(c.Request.Context(), req, nil)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AIHandler) ChatStream(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Question == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入问题")
		return
	}
	flusher, ok := sseWriter(c)
	if !ok {
		return
	}
	result, err := h.svc.Chat(c.Request.Context(), req, func(token string) {
		sseSend(c, flusher, gin.H{"type": "token", "content": token})
	})
	if err != nil {
		sseSend(c, flusher, gin.H{"type": "error", "message": err.Error()})
		return
	}
	sseSend(c, flusher, gin.H{
		"type":           "done",
		"text":           result.Text,
		"sources":        result.Sources,
		"model":          result.Model,
		"provider":       result.Provider,
		"virtualModel":   result.VirtualModel,
		"canonicalModel": result.CanonicalModel,
	})
}

func (h *AIHandler) CareerInterview(c *gin.Context) {
	var req model.CareerInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Message == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入内容")
		return
	}
	flusher, ok := sseWriter(c)
	if !ok {
		return
	}
	text, err := h.svc.CareerInterview(c.Request.Context(), req.Message, req.History, func(token string) {
		sseSend(c, flusher, gin.H{"type": "token", "content": token})
	})
	if err != nil {
		sseSend(c, flusher, gin.H{"type": "error", "message": err.Error()})
		return
	}
	sseSend(c, flusher, gin.H{"type": "done", "text": text})
}

func (h *AIHandler) CareerRecommend(c *gin.Context) {
	var req model.CareerRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	result, err := h.svc.CareerRecommend(c.Request.Context(), req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AIHandler) GoalDecompose(c *gin.Context) {
	var req model.GoalDecomposeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	result, err := h.svc.GoalDecompose(c.Request.Context(), req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AIHandler) LearningSuggest(c *gin.Context) {
	var req model.LearningSuggestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	result, err := h.svc.LearningSuggest(c.Request.Context(), req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, result)
}

// sseWriter 设置 SSE 响应头并返回 flusher。
func sseWriter(c *gin.Context) (http.Flusher, bool) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.Fail(c, http.StatusInternalServerError, 500, "不支持流式输出")
		return nil, false
	}
	return flusher, true
}

func sseSend(c *gin.Context, flusher http.Flusher, payload gin.H) {
	data, _ := json.Marshal(payload)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()
}
