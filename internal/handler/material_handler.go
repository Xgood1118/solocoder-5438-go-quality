package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type MaterialHandler struct{}

func NewMaterialHandler() *MaterialHandler {
	return &MaterialHandler{}
}

func (h *MaterialHandler) List(c *gin.Context) {
	materials := store.GlobalStore.ListMaterials()
	c.JSON(http.StatusOK, materials)
}

func (h *MaterialHandler) Get(c *gin.Context) {
	id := c.Param("id")
	material, ok := store.GlobalStore.GetMaterial(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
		return
	}
	c.JSON(http.StatusOK, material)
}

func (h *MaterialHandler) Create(c *gin.Context) {
	var material model.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material.ID = uuid.New().String()
	material.CreatedAt = time.Now()
	material.UpdatedAt = time.Now()

	if material.Type == "" {
		material.Type = model.MaterialMetal
	}

	store.GlobalStore.SaveMaterial(&material)

	h.createDefaultStandard(material.ID, material.Type)

	c.JSON(http.StatusCreated, material)
}

func (h *MaterialHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, ok := store.GlobalStore.GetMaterial(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
		return
	}

	var material model.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material.ID = existing.ID
	material.CreatedAt = existing.CreatedAt
	material.UpdatedAt = time.Now()

	store.GlobalStore.Materials.Store(id, &material)
	c.JSON(http.StatusOK, material)
}

func (h *MaterialHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	store.GlobalStore.DeleteMaterial(id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *MaterialHandler) createDefaultStandard(materialID string, materialType model.MaterialType) {
	defaultItems := h.getDefaultItems(materialType)

	types := []model.InspectionType{model.TypeIQC, model.TypeOQC, model.TypeOBC}
	for _, t := range types {
		std := &model.Standard{
			ID:               uuid.New().String(),
			MaterialID:       materialID,
			Type:             t,
			InspectionLevels: defaultItems,
			AQL:              1.0,
			InspectionLevel:  "II",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		store.GlobalStore.SaveStandard(std)
	}
}

func (h *MaterialHandler) getDefaultItems(materialType model.MaterialType) []model.InspectionItem {
	items := []model.InspectionItem{
		{ID: uuid.New().String(), Name: "外观", Standard: "无明显划痕、变形", IsNumeric: false, SPCEnabled: false},
		{ID: uuid.New().String(), Name: "尺寸", Standard: "10.0±0.1mm", Unit: "mm", MinValue: 9.9, MaxValue: 10.1, IsNumeric: true, SPCEnabled: true},
		{ID: uuid.New().String(), Name: "性能", Standard: "符合规格要求", IsNumeric: false, SPCEnabled: false},
	}

	switch materialType {
	case model.MaterialMetal:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "硬度", Standard: "HRC 20-30", Unit: "HRC", MinValue: 20, MaxValue: 30, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "化学成分", Standard: "符合标准", IsNumeric: false, SPCEnabled: false},
		)
	case model.MaterialPlastic:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "颜色", Standard: "无色差", IsNumeric: false, SPCEnabled: false},
			model.InspectionItem{ID: uuid.New().String(), Name: "强度", Standard: "≥50MPa", Unit: "MPa", MinValue: 50, MaxValue: 100, IsNumeric: true, SPCEnabled: true},
		)
	case model.MaterialElectronic:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "电阻", Standard: "100±5Ω", Unit: "Ω", MinValue: 95, MaxValue: 105, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "耐压", Standard: "≥500V", Unit: "V", MinValue: 500, MaxValue: 1000, IsNumeric: true, SPCEnabled: false},
		)
	case model.MaterialChemical:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "纯度", Standard: "≥99.5%", Unit: "%", MinValue: 99.5, MaxValue: 100, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "水分", Standard: "≤0.1%", Unit: "%", MinValue: 0, MaxValue: 0.1, IsNumeric: true, SPCEnabled: false},
		)
	case model.MaterialPaper:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "克重", Standard: "200±5g/m²", Unit: "g/m²", MinValue: 195, MaxValue: 205, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "厚度", Standard: "0.2±0.02mm", Unit: "mm", MinValue: 0.18, MaxValue: 0.22, IsNumeric: true, SPCEnabled: false},
		)
	}

	return items
}
