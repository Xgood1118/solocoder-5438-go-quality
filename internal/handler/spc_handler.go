package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qc-system/internal/service"
	"qc-system/internal/store"
)

type SPCHandler struct {
	spcService *service.SPCService
}

func NewSPCHandler() *SPCHandler {
	return &SPCHandler{
		spcService: service.NewSPCService(),
	}
}

func (h *SPCHandler) GetChart(c *gin.Context) {
	materialID := c.Param("material_id")
	itemID := c.Query("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item_id is required"})
		return
	}

	chart, ok := h.spcService.GetSPCChart(materialID, itemID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "chart not found"})
		return
	}
	c.JSON(http.StatusOK, chart)
}

func (h *SPCHandler) UpdateChart(c *gin.Context) {
	materialID := c.Param("material_id")
	itemID := c.Query("item_id")
	itemName := c.Query("item_name")

	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item_id is required"})
		return
	}

	chart, err := h.spcService.UpdateSPCChart(materialID, itemID, itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if chart == nil {
		c.JSON(http.StatusOK, gin.H{"message": "not enough data"})
		return
	}
	c.JSON(http.StatusOK, chart)
}

func (h *SPCHandler) ListCharts(c *gin.Context) {
	charts := store.GlobalStore.ListSPCCharts()
	c.JSON(http.StatusOK, charts)
}

func (h *SPCHandler) CheckAlarms(c *gin.Context) {
	materialID := c.Param("material_id")
	itemID := c.Query("item_id")

	chart, ok := h.spcService.GetSPCChart(materialID, itemID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "chart not found"})
		return
	}

	alarms := h.spcService.CheckAlarms(chart)
	c.JSON(http.StatusOK, gin.H{"alarms": alarms, "count": len(alarms)})
}

func (h *SPCHandler) CalculateControlLimits(c *gin.Context) {
	var req struct {
		Values []float64 `json:"values" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cl, ucl, lcl := h.spcService.CalculateControlLimits(req.Values)
	c.JSON(http.StatusOK, gin.H{
		"cl":  cl,
		"ucl": ucl,
		"lcl": lcl,
	})
}

func (h *SPCHandler) CalculateXBarR(c *gin.Context) {
	var req struct {
		Subgroups [][]float64 `json:"subgroups" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	xBarCL, xBarUCL, xBarLCL, rCL, rUCL, rLCL := h.spcService.CalculateXBarR(req.Subgroups)
	c.JSON(http.StatusOK, gin.H{
		"xbar_cl":  xBarCL,
		"xbar_ucl": xBarUCL,
		"xbar_lcl": xBarLCL,
		"r_cl":     rCL,
		"r_ucl":    rUCL,
		"r_lcl":    rLCL,
	})
}
