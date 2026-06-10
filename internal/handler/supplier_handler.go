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

type SupplierHandler struct {
	scoreService *service.SupplierScoreService
}

func NewSupplierHandler() *SupplierHandler {
	return &SupplierHandler{
		scoreService: service.NewSupplierScoreService(),
	}
}

func (h *SupplierHandler) List(c *gin.Context) {
	suppliers := store.GlobalStore.ListSuppliers()
	c.JSON(http.StatusOK, suppliers)
}

func (h *SupplierHandler) Get(c *gin.Context) {
	id := c.Param("id")
	supplier, ok := store.GlobalStore.GetSupplier(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}
	c.JSON(http.StatusOK, supplier)
}

func (h *SupplierHandler) Create(c *gin.Context) {
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier.ID = uuid.New().String()
	supplier.Status = "normal"
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()
	supplier.Rating = 80

	store.GlobalStore.SaveSupplier(&supplier)
	c.JSON(http.StatusCreated, supplier)
}

func (h *SupplierHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetSupplier(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}

	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier.ID = existing.ID
	supplier.CreatedAt = existing.CreatedAt
	supplier.UpdatedAt = time.Now()

	store.GlobalStore.Suppliers.Store(id, &supplier)
	c.JSON(http.StatusOK, supplier)
}

func (h *SupplierHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	store.GlobalStore.Suppliers.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *SupplierHandler) GetScores(c *gin.Context) {
	id := c.Param("id")
	scores := h.scoreService.GetSupplierScores(id)
	c.JSON(http.StatusOK, scores)
}

func (h *SupplierHandler) CalculateScore(c *gin.Context) {
	id := c.Param("id")
	yearMonth := c.Query("year_month")
	if yearMonth == "" {
		yearMonth = time.Now().Format("200601")
	}

	score, err := h.scoreService.CalculateMonthlyScore(id, yearMonth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, score)
}

func (h *SupplierHandler) GetFocusSuppliers(c *gin.Context) {
	suppliers := h.scoreService.GetFocusSuppliers()
	c.JSON(http.StatusOK, suppliers)
}

func (h *SupplierHandler) CalculateAll(c *gin.Context) {
	yearMonth := c.Query("year_month")
	if yearMonth == "" {
		yearMonth = time.Now().Format("200601")
	}

	scores := h.scoreService.CalculateAllSuppliers(yearMonth)
	c.JSON(http.StatusOK, scores)
}
