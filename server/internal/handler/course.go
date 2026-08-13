package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct{ svc *service.CourseService }

func NewCourseHandler(svc *service.CourseService) *CourseHandler { return &CourseHandler{svc: svc} }

func (h *CourseHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CourseHandler) Get(c *gin.Context) {
	course, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, course)
}

func (h *CourseHandler) Create(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Create(course)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *CourseHandler) Update(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Update(c.Param("id"), course)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *CourseHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

func (h *CourseHandler) AddChapter(c *gin.Context) {
	var ch model.Chapter
	if err := c.ShouldBindJSON(&ch); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.AddChapter(c.Param("id"), ch)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "章节创建成功")
}

func (h *CourseHandler) UpdateChapter(c *gin.Context) {
	var ch model.Chapter
	if err := c.ShouldBindJSON(&ch); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateChapter(c.Param("id"), c.Param("chapterId"), ch)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "章节更新成功")
}

func (h *CourseHandler) DeleteChapter(c *gin.Context) {
	if err := h.svc.DeleteChapter(c.Param("id"), c.Param("chapterId")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "章节删除成功")
}
