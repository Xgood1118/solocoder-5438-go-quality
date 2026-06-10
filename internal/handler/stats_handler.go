package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats := store.GlobalStore.GetStats()
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetQualityStats(c *gin.Context) {
	materialType := c.Query("material_type")
	supplierID := c.Query("supplier_id")
	inspectorID := c.Query("inspector_id")
	typ := model.InspectionType(c.Query("type"))

	records := store.GlobalStore.ListRecords()

	total := 0
	pass := 0
	fail := 0
	reject := 0

	for _, rec := range records {
		if materialType != "" {
			material, ok := store.GlobalStore.GetMaterial(rec.MaterialID)
			if !ok || string(material.Type) != materialType {
				continue
			}
		}
		if supplierID != "" && rec.SupplierID != supplierID {
			continue
		}
		if inspectorID != "" && rec.InspectorID != inspectorID {
			continue
		}
		if typ != "" && rec.Type != typ {
			continue
		}

		total++
		switch rec.FinalJudgment {
		case model.JudgmentPass:
			pass++
		case model.JudgmentFail:
			fail++
		case model.JudgmentReject:
			reject++
		}
	}

	passRate := 0.0
	if total > 0 {
		passRate = float64(pass) / float64(total) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"pass":       pass,
		"fail":       fail,
		"reject":     reject,
		"pass_rate":  passRate,
		"filter": gin.H{
			"material_type": materialType,
			"supplier_id":   supplierID,
			"inspector_id":  inspectorID,
			"type":          typ,
		},
	})
}

func (h *StatsHandler) ByMaterialType(c *gin.Context) {
	type materialStats struct {
		MaterialCount int `json:"material_count"`
		RecordCount   int `json:"record_count"`
		PassCount     int `json:"pass_count"`
	}

	typeStats := make(map[string]*materialStats)

	materials := store.GlobalStore.ListMaterials()
	for _, m := range materials {
		key := string(m.Type)
		if _, ok := typeStats[key]; !ok {
			typeStats[key] = &materialStats{}
		}
		typeStats[key].MaterialCount++
	}

	records := store.GlobalStore.ListRecords()
	for _, rec := range records {
		material, ok := store.GlobalStore.GetMaterial(rec.MaterialID)
		if !ok {
			continue
		}
		key := string(material.Type)
		if _, ok := typeStats[key]; !ok {
			continue
		}
		typeStats[key].RecordCount++
		if rec.FinalJudgment == model.JudgmentPass {
			typeStats[key].PassCount++
		}
	}

	c.JSON(http.StatusOK, typeStats)
}

func (h *StatsHandler) BySupplier(c *gin.Context) {
	suppliers := store.GlobalStore.ListSuppliers()
	result := make([]gin.H, 0, len(suppliers))

	for _, sup := range suppliers {
		records := store.GlobalStore.ListRecords()
		total := 0
		pass := 0

		for _, rec := range records {
			if rec.SupplierID == sup.ID {
				total++
				if rec.FinalJudgment == model.JudgmentPass {
					pass++
				}
			}
		}

		passRate := 0.0
		if total > 0 {
			passRate = float64(pass) / float64(total) * 100
		}

		result = append(result, gin.H{
			"supplier_id":   sup.ID,
			"supplier_name": sup.Name,
			"total":         total,
			"pass":          pass,
			"pass_rate":     passRate,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *StatsHandler) ByInspector(c *gin.Context) {
	inspectors := store.GlobalStore.ListInspectors()
	result := make([]gin.H, 0, len(inspectors))

	for _, ins := range inspectors {
		records := store.GlobalStore.ListRecords()
		total := 0
		pass := 0

		for _, rec := range records {
			if rec.InspectorID == ins.ID {
				total++
				if rec.FinalJudgment == model.JudgmentPass {
					pass++
				}
			}
		}

		result = append(result, gin.H{
			"inspector_id":   ins.ID,
			"inspector_name": ins.Name,
			"total":         total,
			"pass":          pass,
		})
	}

	c.JSON(http.StatusOK, result)
}
