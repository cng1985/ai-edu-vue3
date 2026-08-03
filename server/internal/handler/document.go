package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct{ svc *service.DocumentService }

func NewDocumentHandler(svc *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("type"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *DocumentHandler) Get(c *gin.Context) {
	doc, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, doc)
}

func (h *DocumentHandler) Create(c *gin.Context) {
	var req struct {
		DocNo  string  `json:"docNo"`
		Title  string  `json:"title"`
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
		Status string  `json:"status"`
		Remark string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	claims := middleware.GetClaims(c)
	doc := &model.Document{
		DocNo: req.DocNo, Title: req.Title, Type: req.Type,
		Amount: req.Amount, Status: req.Status, Remark: req.Remark,
		CreatedBy: claims.Username,
	}
	result, err := h.svc.Create(doc)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, result, "创建成功")
}

func (h *DocumentHandler) Update(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	doc, err := h.svc.Update(c.Param("id"), updates)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, doc, "更新成功")
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

func (h *DocumentHandler) Export(c *gin.Context) {
	data, filename, err := h.svc.ExportExcel(c.Query("keyword"), c.Query("type"), c.Query("status"))
	if err != nil {
		failErr(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *DocumentHandler) ExportTemplate(c *gin.Context) {
	data, err := h.svc.ExportTemplate()
	if err != nil {
		failErr(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape("单据导入模板.xlsx"))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *DocumentHandler) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "请上传 Excel 文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "无法读取文件")
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "无法读取文件")
		return
	}

	claims := middleware.GetClaims(c)
	taskID, err := h.svc.StartImport(data, claims.Username)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, gin.H{"taskId": taskID}, "导入任务已启动")
}

func (h *DocumentHandler) ImportProgress(c *gin.Context) {
	progress, err := h.svc.GetImportProgress(c.Param("taskId"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, progress)
}
