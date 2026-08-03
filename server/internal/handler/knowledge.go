package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type KnowledgeHandler struct{ svc *service.KnowledgeService }

func NewKnowledgeHandler(svc *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{svc: svc}
}

func (h *KnowledgeHandler) Status(c *gin.Context) {
	status, err := h.svc.Status()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, status)
}

func (h *KnowledgeHandler) ListChunks(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListChunks(page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *KnowledgeHandler) Reindex(c *gin.Context) {
	status, err := h.svc.RebuildIndex(c.Request.Context())
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, status, "知识库索引重建完成")
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入检索关键词")
		return
	}
	results, err := h.svc.Search(c.Request.Context(), query, 5)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, results)
}
