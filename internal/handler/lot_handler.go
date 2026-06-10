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

type LotHandler struct {
	lotGen *service.LotGenerator
}

func NewLotHandler() *LotHandler {
	return &LotHandler{
		lotGen: service.GetLotGenerator(),
	}
}

func (h *LotHandler) List(c *gin.Context) {
	lots := store.GlobalStore.ListLots()
	c.JSON(http.StatusOK, lots)
}

func (h *LotHandler) Get(c *gin.Context) {
	id := c.Param("id")
	lot, ok := store.GlobalStore.GetLot(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot not found"})
		return
	}
	c.JSON(http.StatusOK, lot)
}

func (h *LotHandler) Create(c *gin.Context) {
	var lot model.Lot
	if err := c.ShouldBindJSON(&lot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material, ok := store.GlobalStore.GetMaterial(lot.MaterialID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "material not found"})
		return
	}

	lot.ID = uuid.New().String()
	lot.LotNo = h.lotGen.GenerateLotNo(material.Code)
	lot.Status = "pending"
	lot.CreatedAt = time.Now()

	if lot.ProductionDate.IsZero() {
		lot.ProductionDate = time.Now()
	}

	store.GlobalStore.SaveLot(&lot)
	c.JSON(http.StatusCreated, lot)
}

func (h *LotHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetLot(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot not found"})
		return
	}

	var lot model.Lot
	if err := c.ShouldBindJSON(&lot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lot.ID = existing.ID
	lot.LotNo = existing.LotNo
	lot.CreatedAt = existing.CreatedAt

	store.GlobalStore.Lots.Store(id, &lot)
	c.JSON(http.StatusOK, lot)
}

func (h *LotHandler) GetByLotNo(c *gin.Context) {
	lotNo := c.Param("lot_no")
	lot, ok := store.GlobalStore.GetLotByNo(lotNo)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot not found"})
		return
	}
	c.JSON(http.StatusOK, lot)
}

func (h *LotHandler) GetRecords(c *gin.Context) {
	id := c.Param("id")
	records := store.GlobalStore.GetRecordsByLot(id)
	c.JSON(http.StatusOK, records)
}
