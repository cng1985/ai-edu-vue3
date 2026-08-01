package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AIHandler struct{ svc *service.AIService }
type RBACHandler struct{ svc *service.RBACService }
type AppHandler struct {
	courses *service.CourseService
	quizzes *service.QuizService
}

func NewAIHandler(svc *service.AIService) *AIHandler       { return &AIHandler{svc: svc} }
func NewRBACHandler(svc *service.RBACService) *RBACHandler { return &RBACHandler{svc: svc} }
func NewAppHandler(courses *service.CourseService, quizzes *service.QuizService) *AppHandler {
	return &AppHandler{courses: courses, quizzes: quizzes}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Register(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "注册成功")
}

func (h *AuthHandler) GuestLogin(c *gin.Context) {
	res, err := h.svc.GuestLogin()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AIHandler) Config(c *gin.Context) {
	response.OK(c, h.svc.ConfigInfo())
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Question == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入问题")
		return
	}
	result, err := h.svc.Chat(c.Request.Context(), req.Question, req.History, nil)
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
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.Fail(c, http.StatusInternalServerError, 500, "不支持流式输出")
		return
	}
	result, err := h.svc.Chat(c.Request.Context(), req.Question, req.History, func(token string) {
		data, _ := json.Marshal(gin.H{"type": "token", "content": token})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	})
	if err != nil {
		data, _ := json.Marshal(gin.H{"type": "error", "message": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
		return
	}
	done, _ := json.Marshal(gin.H{"type": "done", "text": result.Text, "sources": result.Sources})
	fmt.Fprintf(c.Writer, "data: %s\n\n", done)
	flusher.Flush()
}

func (h *AIHandler) CareerInterview(c *gin.Context) {
	var req model.CareerInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Message == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入内容")
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.Fail(c, http.StatusInternalServerError, 500, "不支持流式输出")
		return
	}
	text, err := h.svc.CareerInterview(c.Request.Context(), req.Message, req.History, func(token string) {
		data, _ := json.Marshal(gin.H{"type": "token", "content": token})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	})
	if err != nil {
		data, _ := json.Marshal(gin.H{"type": "error", "message": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
		return
	}
	done, _ := json.Marshal(gin.H{"type": "done", "text": text})
	fmt.Fprintf(c.Writer, "data: %s\n\n", done)
	flusher.Flush()
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

func (h *RBACHandler) ListPermissions(c *gin.Context) {
	response.OK(c, h.svc.ListPermissions())
}

func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.svc.ListRoles()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, roles)
}

func (h *RBACHandler) UpdateRole(c *gin.Context) {
	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateRole(c.Param("role"), req.Permissions); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "权限已更新")
}

func (h *AppHandler) ListCourses(c *gin.Context) {
	courses, err := h.courses.ListPublished()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, courses)
}

func (h *AppHandler) GetCourse(c *gin.Context) {
	course, err := h.courses.GetPublished(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, course)
}

func (h *AppHandler) GetQuiz(c *gin.Context) {
	quiz, err := h.quizzes.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	if quiz.Status != "published" {
		response.Fail(c, http.StatusNotFound, 404, "测验不存在")
		return
	}
	response.OK(c, quiz)
}

func (h *AppHandler) Profile(c *gin.Context) {
	claims := middleware.GetClaims(c)
	response.OK(c, gin.H{
		"id": claims.ID, "username": claims.Username, "role": claims.Role,
	})
}
