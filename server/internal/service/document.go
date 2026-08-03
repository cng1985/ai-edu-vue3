package service

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/xuri/excelize/v2"
)

type DocumentService struct {
	docs *repository.DocumentRepo
}

func NewDocumentService(docs *repository.DocumentRepo) *DocumentService {
	return &DocumentService{docs: docs}
}

var (
	importTasks   = sync.Map{}
	importTaskTTL = 30 * time.Minute
)

func (s *DocumentService) List(keyword, docType, status string, page, pageSize int) (*model.PageResult[model.Document], error) {
	list, total, err := s.docs.List(keyword, docType, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.Document]{List: list, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *DocumentService) Get(id string) (*model.Document, error) {
	doc, err := s.docs.FindByID(id)
	if err != nil {
		return nil, errors.New("单据不存在")
	}
	return doc, nil
}

func (s *DocumentService) Create(doc *model.Document) (*model.Document, error) {
	if doc.DocNo == "" || doc.Title == "" {
		return nil, errors.New("单据编号和标题不能为空")
	}
	existing, _ := s.docs.FindByDocNo(doc.DocNo)
	if existing != nil && existing.ID != "" {
		return nil, errors.New("单据编号已存在")
	}
	now := time.Now().UnixMilli()
	doc.ID = genID("doc")
	if doc.Status == "" {
		doc.Status = "draft"
	}
	doc.CreatedAt = now
	doc.UpdatedAt = now
	if err := s.docs.Create(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *DocumentService) Update(id string, updates map[string]interface{}) (*model.Document, error) {
	doc, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if v, ok := updates["docNo"].(string); ok && v != "" {
		existing, _ := s.docs.FindByDocNo(v)
		if existing != nil && existing.ID != id {
			return nil, errors.New("单据编号已存在")
		}
		doc.DocNo = v
	}
	if v, ok := updates["title"].(string); ok {
		doc.Title = v
	}
	if v, ok := updates["type"].(string); ok {
		doc.Type = v
	}
	if v, ok := updates["amount"].(float64); ok {
		doc.Amount = v
	}
	if v, ok := updates["status"].(string); ok {
		doc.Status = v
	}
	if v, ok := updates["remark"].(string); ok {
		doc.Remark = v
	}
	doc.UpdatedAt = time.Now().UnixMilli()
	if err := s.docs.Update(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *DocumentService) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	return s.docs.Delete(id)
}

func (s *DocumentService) ExportExcel(keyword, docType, status string) ([]byte, string, error) {
	docs, err := s.docs.ListAll(keyword, docType, status)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	sheet := "单据列表"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"单据编号", "标题", "类型", "金额", "状态", "备注", "创建人", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	statusMap := map[string]string{
		"draft": "草稿", "pending": "待审核", "approved": "已通过", "rejected": "已驳回",
	}
	typeMap := map[string]string{
		"purchase": "采购单", "sales": "销售单", "expense": "报销单", "other": "其他",
	}

	for i, doc := range docs {
		row := i + 2
		statusLabel := statusMap[doc.Status]
		if statusLabel == "" {
			statusLabel = doc.Status
		}
		typeLabel := typeMap[doc.Type]
		if typeLabel == "" {
			typeLabel = doc.Type
		}
		createdAt := time.UnixMilli(doc.CreatedAt).Format("2006-01-02 15:04:05")

		values := []interface{}{doc.DocNo, doc.Title, typeLabel, doc.Amount, statusLabel, doc.Remark, doc.CreatedBy, createdAt}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// 设置列宽
	colWidths := []float64{18, 30, 12, 12, 10, 30, 12, 20}
	for i, w := range colWidths {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, w)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("单据导出_%s.xlsx", time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}

func (s *DocumentService) StartImport(fileData []byte, createdBy string) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return "", errors.New("无法解析 Excel 文件")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return "", errors.New("Excel 文件为空或缺少数据行")
	}

	taskID := genID("import")
	progress := &model.ImportTaskProgress{
		TaskID:  taskID,
		Status:  "pending",
		Total:   len(rows) - 1,
		Current: 0,
	}
	importTasks.Store(taskID, progress)

	go s.processImport(taskID, rows[1:], createdBy)
	return taskID, nil
}

func (s *DocumentService) processImport(taskID string, rows [][]string, createdBy string) {
	val, _ := importTasks.Load(taskID)
	progress := val.(*model.ImportTaskProgress)
	progress.Status = "processing"

	typeMap := map[string]string{
		"采购单": "purchase", "销售单": "sales", "报销单": "expense", "其他": "other",
	}
	statusMap := map[string]string{
		"草稿": "draft", "待审核": "pending", "已通过": "approved", "已驳回": "rejected",
	}

	for i, row := range rows {
		progress.Current = i + 1

		if len(row) < 2 || strings.TrimSpace(row[0]) == "" {
			progress.Failed++
			progress.Errors = append(progress.Errors, fmt.Sprintf("第 %d 行：单据编号不能为空", i+2))
			continue
		}

		docNo := strings.TrimSpace(row[0])
		title := ""
		if len(row) > 1 {
			title = strings.TrimSpace(row[1])
		}
		if title == "" {
			progress.Failed++
			progress.Errors = append(progress.Errors, fmt.Sprintf("第 %d 行：标题不能为空", i+2))
			continue
		}

		docType := "other"
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			t := strings.TrimSpace(row[2])
			if v, ok := typeMap[t]; ok {
				docType = v
			} else {
				docType = t
			}
		}

		var amount float64
		if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
			a, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
			if err != nil {
				progress.Failed++
				progress.Errors = append(progress.Errors, fmt.Sprintf("第 %d 行：金额格式错误", i+2))
				continue
			}
			amount = a
		}

		status := "draft"
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			st := strings.TrimSpace(row[4])
			if v, ok := statusMap[st]; ok {
				status = v
			} else {
				status = st
			}
		}

		remark := ""
		if len(row) > 5 {
			remark = strings.TrimSpace(row[5])
		}

		existing, _ := s.docs.FindByDocNo(docNo)
		now := time.Now().UnixMilli()
		if existing != nil && existing.ID != "" {
			existing.Title = title
			existing.Type = docType
			existing.Amount = amount
			existing.Status = status
			existing.Remark = remark
			existing.UpdatedAt = now
			if err := s.docs.Update(existing); err != nil {
				progress.Failed++
				progress.Errors = append(progress.Errors, fmt.Sprintf("第 %d 行：更新失败 - %s", i+2, err.Error()))
				continue
			}
		} else {
			doc := &model.Document{
				ID: genID("doc"), DocNo: docNo, Title: title, Type: docType,
				Amount: amount, Status: status, Remark: remark,
				CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now,
			}
			if err := s.docs.Create(doc); err != nil {
				progress.Failed++
				progress.Errors = append(progress.Errors, fmt.Sprintf("第 %d 行：创建失败 - %s", i+2, err.Error()))
				continue
			}
		}
		progress.Success++

		// 模拟处理延迟，使进度条可见
		time.Sleep(50 * time.Millisecond)
	}

	progress.Status = "completed"
	progress.Message = fmt.Sprintf("导入完成：成功 %d 条，失败 %d 条", progress.Success, progress.Failed)

	go func() {
		time.Sleep(importTaskTTL)
		importTasks.Delete(taskID)
	}()
}

func (s *DocumentService) GetImportProgress(taskID string) (*model.ImportTaskProgress, error) {
	val, ok := importTasks.Load(taskID)
	if !ok {
		return nil, errors.New("导入任务不存在或已过期")
	}
	return val.(*model.ImportTaskProgress), nil
}

func (s *DocumentService) ExportTemplate() ([]byte, error) {
	f := excelize.NewFile()
	sheet := "导入模板"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"单据编号", "标题", "类型", "金额", "状态", "备注"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	examples := [][]interface{}{
		{"DOC-2026-001", "办公用品采购", "采购单", 1500.00, "草稿", "季度办公用品"},
		{"DOC-2026-002", "软件授权销售", "销售单", 9800.00, "待审核", ""},
		{"DOC-2026-003", "差旅报销", "报销单", 2350.50, "已通过", "北京出差"},
	}
	for i, row := range examples {
		for j, v := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	colWidths := []float64{18, 30, 12, 12, 10, 30}
	for i, w := range colWidths {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, w)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
