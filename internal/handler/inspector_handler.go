package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/service"
	"qc-system/internal/store"
)

type InspectorHandler struct {
	delegationService *service.DelegationService
	recordService     *service.RecordService
}

func NewInspectorHandler() *InspectorHandler {
	return &InspectorHandler{
		delegationService: service.NewDelegationService(),
		recordService:     service.NewRecordService(),
	}
}

func (h *InspectorHandler) List(c *gin.Context) {
	inspectors := store.GlobalStore.ListInspectors()
	c.JSON(http.StatusOK, inspectors)
}

func (h *InspectorHandler) Get(c *gin.Context) {
	id := c.Param("id")
	inspector, ok := store.GlobalStore.GetInspector(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "inspector not found"})
		return
	}
	c.JSON(http.StatusOK, inspector)
}

func (h *InspectorHandler) Create(c *gin.Context) {
	var inspector model.Inspector
	if err := c.ShouldBindJSON(&inspector); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inspector.ID = uuid.New().String()
	inspector.Status = "active"
	inspector.CreatedAt = time.Now()
	inspector.UpdatedAt = time.Now()

	store.GlobalStore.SaveInspector(&inspector)
	c.JSON(http.StatusCreated, inspector)
}

func (h *InspectorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetInspector(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "inspector not found"})
		return
	}

	var inspector model.Inspector
	if err := c.ShouldBindJSON(&inspector); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inspector.ID = existing.ID
	inspector.CreatedAt = existing.CreatedAt
	inspector.UpdatedAt = time.Now()

	store.GlobalStore.Inspectors.Store(id, &inspector)
	c.JSON(http.StatusOK, inspector)
}

func (h *InspectorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	store.GlobalStore.Inspectors.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *InspectorHandler) GetRecords(c *gin.Context) {
	id := c.Param("id")
	records := store.GlobalStore.GetRecordsByInspector(id)
	c.JSON(http.StatusOK, records)
}

func (h *InspectorHandler) CreateDelegation(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		DelegateeID string `json:"delegatee_id" binding:"required"`
		Days        int    `json:"days" binding:"required"`
		Role        string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "inspector"
	}

	delegation, err := h.delegationService.CreateDelegation(id, req.DelegateeID, req.Role, req.Days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, delegation)
}

func (h *InspectorHandler) GetActiveDelegation(c *gin.Context) {
	id := c.Param("id")
	delegation, ok := h.delegationService.GetActiveDelegation(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active delegation"})
		return
	}
	c.JSON(http.StatusOK, delegation)
}

func (h *InspectorHandler) RevokeDelegation(c *gin.Context) {
	id := c.Param("id")
	err := h.delegationService.RevokeDelegation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "delegation revoked"})
}

func (h *InspectorHandler) TransferRecords(c *gin.Context) {
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
