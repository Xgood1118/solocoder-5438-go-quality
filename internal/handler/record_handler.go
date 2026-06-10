package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"

	"qc-system/internal/model"
	"qc-system/internal/service"
	"qc-system/internal/store"
)

type RecordHandler struct {
	recordService *service.RecordService
	pdfService    *service.PDFService
}

func NewRecordHandler() *RecordHandler {
	return &RecordHandler{
		recordService: service.NewRecordService(),
		pdfService:    service.NewPDFService(),
	}
}

func (h *RecordHandler) List(c *gin.Context) {
	records := store.GlobalStore.ListRecords()
	c.JSON(http.StatusOK, records)
}

func (h *RecordHandler) Get(c *gin.Context) {
	id := c.Param("id")
	record, ok := store.GlobalStore.GetRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Create(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lot, ok := store.GlobalStore.GetLot(req.LotID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lot not found"})
		return
	}

	standard, ok := store.GlobalStore.GetStandardByMaterial(lot.MaterialID, req.Type)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "standard not found"})
		return
	}

	sampleSize := len(standard.InspectionLevels)
	if req.Type == model.TypeOQC && lot.Quantity > 1000 {
		sampleSize = lot.Quantity / 10
	}

	items := make([]model.InspectionItemResult, 0, len(standard.InspectionLevels))
	for _, item := range standard.InspectionLevels {
		items = append(items, model.InspectionItemResult{
			ItemID:    item.ID,
			ItemName:  item.Name,
			Standard:  item.Standard,
			IsNumeric: item.IsNumeric,
			Judgment:  "",
		})
	}

	record := &model.InspectionRecord{
		ID:               uuid.New().String(),
		LotID:            req.LotID,
		LotNo:            lot.LotNo,
		MaterialID:       lot.MaterialID,
		SupplierID:       lot.SupplierID,
		Type:             req.Type,
		IPQCType:         req.IPQCType,
		InspectorID:      req.InspectorID,
		InspectorName:    req.InspectorName,
		Status:           model.StatusDraft,
		TotalSampleSize:  sampleSize,
		DefectCount:      0,
		FinalJudgment:    "",
		Items:            items,
		Version:          1,
		Remarks:          []model.Remark{},
		IsFullInspection: req.Type == model.TypeOQC && lot.Quantity <= 1000,
		CreatedAt:        time.Now(),
	}

	store.GlobalStore.SaveRecord(record)
	c.JSON(http.StatusCreated, record)
}

type CreateRecordRequest struct {
	LotID         string              `json:"lot_id" binding:"required"`
	Type          model.InspectionType `json:"type" binding:"required"`
	IPQCType      model.IPQCType      `json:"ipqc_type"`
	InspectorID   string              `json:"inspector_id" binding:"required"`
	InspectorName string              `json:"inspector_name"`
}

func (h *RecordHandler) Update(c *gin.Context) {
	id := c.Param("id")
	record, ok := store.GlobalStore.GetRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	if record.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only draft records can be edited"})
		return
	}

	var updateData struct {
		Items         []model.InspectionItemResult `json:"items"`
		DefectCount   int                          `json:"defect_count"`
		FinalJudgment model.JudgmentResult         `json:"final_judgment"`
		Disposition   model.DispositionType        `json:"disposition"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Items != nil {
		record.Items = updateData.Items
	}
	if updateData.DefectCount > 0 {
		record.DefectCount = updateData.DefectCount
	}
	if updateData.FinalJudgment != "" {
		record.FinalJudgment = updateData.FinalJudgment

		lot, _ := store.GlobalStore.GetLot(record.LotID)
		if record.Type == model.TypeOQC && lot != nil && lot.Quantity > 1000 && !record.IsFullInspection {
			if updateData.FinalJudgment == model.JudgmentFail || updateData.FinalJudgment == model.JudgmentReject {
				record.IsFullInspection = true
				record.TotalSampleSize = lot.Quantity

				remark := model.Remark{
					ID:        uuid.New().String(),
					Content:   "抽检不合格，已自动转为全检（原10%抽检→全量检验）",
					Operator:  "system",
					CreatedAt: time.Now(),
				}
				record.Remarks = append(record.Remarks, remark)

				auditLog := &model.AuditLog{
					ID:        uuid.New().String(),
					RecordID:  record.ID,
					Action:    "oqc_sampling_to_full",
					OldStatus: string(record.Status),
					NewStatus: string(record.Status),
					Operator:  "system",
					Remark:    "OQC抽检不合格自动转全检",
					CreatedAt: time.Now(),
				}
				store.GlobalStore.SaveAuditLog(auditLog)
			}
		}
	}
	if updateData.Disposition != "" {
		record.Disposition = updateData.Disposition
	}

	store.GlobalStore.UpdateRecord(record)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Submit(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "system"
	}

	err := h.recordService.SubmitRecord(id, operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	role := c.GetHeader("X-Role")
	if operator == "" {
		operator = "system"
	}
	if role == "" {
		role = "qa_manager"
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	err := h.recordService.ApproveRecord(id, operator, role, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	role := c.GetHeader("X-Role")
	if operator == "" {
		operator = "system"
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	err := h.recordService.RejectRecord(id, operator, role, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Close(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "system"
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	err := h.recordService.CloseRecord(id, operator, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) AddRemark(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "system"
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.recordService.AddRemark(id, req.Content, operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	c.JSON(http.StatusOK, record)
}

func (h *RecordHandler) Reinspect(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "system"
	}

	newRecord, err := h.recordService.CreateReinspection(id, operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newRecord)
}

func (h *RecordHandler) Transfer(c *gin.Context) {
	var req struct {
		FromInspectorID string `json:"from_inspector_id" binding:"required"`
		ToInspectorID   string `json:"to_inspector_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "admin"
	}

	count, err := h.recordService.TransferRecords(req.FromInspectorID, req.ToInspectorID, operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transferred_count": count})
}

func (h *RecordHandler) GetAuditLogs(c *gin.Context) {
	id := c.Param("id")
	logs := store.GlobalStore.GetAuditLogsByRecord(id)
	c.JSON(http.StatusOK, logs)
}

func (h *RecordHandler) ExportPDF(c *gin.Context) {
	id := c.Param("id")

	pdfData, err := h.pdfService.GenerateInspectionReport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	record, _ := store.GlobalStore.GetRecord(id)
	fileName := h.pdfService.GetReportFileName(id, record.LotNo)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

func (h *RecordHandler) CreateConcessionApproval(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetHeader("X-Operator")
	if operator == "" {
		operator = "system"
	}

	approval, err := h.recordService.CreateConcessionApproval(id, operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"approval_id": approval.ID,
		"status":      approval.Status,
		"message":     "让步接收审批已创建，等待质量经理初审",
	})
}

func (h *RecordHandler) GetApproval(c *gin.Context) {
	id := c.Param("id")
	record, ok := store.GlobalStore.GetRecord(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	if record.ApprovalID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "this record has no approval flow"})
		return
	}
	approval, ok := store.GlobalStore.GetApproval(record.ApprovalID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "approval not found"})
		return
	}
	c.JSON(http.StatusOK, approval)
}

func (h *RecordHandler) GetByLot(c *gin.Context) {
	lotID := c.Param("lot_id")
	records := store.GlobalStore.GetRecordsByLot(lotID)
	c.JSON(http.StatusOK, records)
}
