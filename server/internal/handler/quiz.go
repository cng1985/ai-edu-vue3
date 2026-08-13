package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct{ svc *service.QuizService }

func NewQuizHandler(svc *service.QuizService) *QuizHandler { return &QuizHandler{svc: svc} }

func (h *QuizHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("courseId"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *QuizHandler) Get(c *gin.Context) {
	quiz, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, quiz)
}

func (h *QuizHandler) Create(c *gin.Context) {
	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Create(quiz)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *QuizHandler) Update(c *gin.Context) {
	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Update(c.Param("id"), quiz)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *QuizHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}
