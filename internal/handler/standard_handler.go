package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/service/aql"
	"qc-system/internal/store"
)

type StandardHandler struct{}

func NewStandardHandler() *StandardHandler {
	return &StandardHandler{}
}

func (h *StandardHandler) List(c *gin.Context) {
	standards := store.GlobalStore.ListStandards()
	c.JSON(http.StatusOK, standards)
}

func (h *StandardHandler) Get(c *gin.Context) {
	id := c.Param("id")
	standard, ok := store.GlobalStore.GetStandard(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
		return
	}
	c.JSON(http.StatusOK, standard)
}

func (h *StandardHandler) Create(c *gin.Context) {
	var standard model.Standard
	if err := c.ShouldBindJSON(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	standard.ID = uuid.New().String()
	standard.CreatedAt = time.Now()
	standard.UpdatedAt = time.Now()

	for i := range standard.InspectionLevels {
		if standard.InspectionLevels[i].ID == "" {
			standard.InspectionLevels[i].ID = uuid.New().String()
		}
	}

	store.GlobalStore.SaveStandard(&standard)
	c.JSON(http.StatusCreated, standard)
}

func (h *StandardHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetStandard(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
		return
	}

	var standard model.Standard
	if err := c.ShouldBindJSON(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	standard.ID = existing.ID
	standard.MaterialID = existing.MaterialID
	standard.Type = existing.Type
	standard.CreatedAt = existing.CreatedAt
	standard.UpdatedAt = time.Now()

	store.GlobalStore.Standards.Store(id, &standard)
	c.JSON(http.StatusOK, standard)
}

func (h *StandardHandler) GetByMaterial(c *gin.Context) {
	materialID := c.Param("material_id")
	typ := model.InspectionType(c.Query("type"))
	if typ == "" {
		typ = model.TypeIQC
	}

	standard, ok := store.GlobalStore.GetStandardByMaterial(materialID, typ)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
		return
	}
	c.JSON(http.StatusOK, standard)
}

func (h *StandardHandler) CalculateSample(c *gin.Context) {
	var req struct {
		BatchSize       int     `json:"batch_size" binding:"required"`
		AQL             float64 `json:"aql" binding:"required"`
		InspectionLevel string  `json:"inspection_level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.InspectionLevel == "" {
		req.InspectionLevel = "II"
	}

	result, err := aql.CalculateSample(req.BatchSize, req.AQL, req.InspectionLevel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sample_size":   result.SampleSize,
		"accept_num":    result.AcceptNum,
		"reject_num":    result.RejectNum,
		"rounded_batch": aql.RoundBatchSize(req.BatchSize),
	})
}

func (h *StandardHandler) GetAQLTable(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"batch_ranges": aql.GetBatchRanges(),
		"aql_values":   aql.GetAQLValues(),
	})
}
