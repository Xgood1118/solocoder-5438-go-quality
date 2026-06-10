package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type NCHandler struct{}

func NewNCHandler() *NCHandler {
	return &NCHandler{}
}

func (h *NCHandler) List(c *gin.Context) {
	ncs := store.GlobalStore.ListNonconformances()
	c.JSON(http.StatusOK, ncs)
}

func (h *NCHandler) Get(c *gin.Context) {
	id := c.Param("id")
	nc, ok := store.GlobalStore.GetNonconformance(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "nonconformance not found"})
		return
	}
	c.JSON(http.StatusOK, nc)
}

func (h *NCHandler) Create(c *gin.Context) {
	var nc model.Nonconformance
	if err := c.ShouldBindJSON(&nc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nc.ID = uuid.New().String()
	if nc.NCRNo == "" {
		nc.NCRNo = "NCR-" + time.Now().Format("200601") + "-" + uuid.New().String()[:6]
	}
	if nc.Status == "" {
		nc.Status = "open"
	}
	nc.CreatedAt = time.Now()

	store.GlobalStore.SaveNonconformance(&nc)
	c.JSON(http.StatusCreated, nc)
}

func (h *NCHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetNonconformance(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "nonconformance not found"})
		return
	}

	var nc model.Nonconformance
	if err := c.ShouldBindJSON(&nc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nc.ID = existing.ID
	nc.NCRNo = existing.NCRNo
	nc.CreatedAt = existing.CreatedAt

	if nc.Status == "closed" && existing.Status != "closed" {
		now := time.Now()
		nc.ClosedAt = &now
	}

	store.GlobalStore.Nonconformances.Store(id, &nc)
	c.JSON(http.StatusOK, nc)
}

func (h *NCHandler) Close(c *gin.Context) {
	id := c.Param("id")
	nc, ok := store.GlobalStore.GetNonconformance(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "nonconformance not found"})
		return
	}

	nc.Status = "closed"
	now := time.Now()
	nc.ClosedAt = &now

	store.GlobalStore.Nonconformances.Store(id, nc)
	c.JSON(http.StatusOK, nc)
}

func (h *NCHandler) ListReturnOrders(c *gin.Context) {
	ros := store.GlobalStore.ListReturnOrders()
	c.JSON(http.StatusOK, ros)
}

func (h *NCHandler) GetReturnOrder(c *gin.Context) {
	id := c.Param("id")
	ro, ok := store.GlobalStore.GetReturnOrder(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "return order not found"})
		return
	}
	c.JSON(http.StatusOK, ro)
}

func (h *NCHandler) CreateReturnOrder(c *gin.Context) {
	var ro model.ReturnOrder
	if err := c.ShouldBindJSON(&ro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ro.ID = uuid.New().String()
	if ro.ReturnNo == "" {
		ro.ReturnNo = "RTN-" + time.Now().Format("200601") + "-" + uuid.New().String()[:6]
	}
	if ro.Status == "" {
		ro.Status = "created"
	}
	ro.CreatedAt = time.Now()

	store.GlobalStore.SaveReturnOrder(&ro)
	c.JSON(http.StatusCreated, ro)
}
