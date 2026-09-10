package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// AppHandler 提供学员端只读接口（已发布的课程与测验）。
type AppHandler struct {
	courses *service.CourseService
	quizzes *service.QuizService
}

func NewAppHandler(courses *service.CourseService, quizzes *service.QuizService) *AppHandler {
	return &AppHandler{courses: courses, quizzes: quizzes}
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
