package model

import "time"

type Material struct {
	ID            string       `json:"id"`
	Code          string       `json:"code"`
	Name          string       `json:"name"`
	Type          MaterialType `json:"type"`
	Specification string       `json:"specification"`
	Unit          string       `json:"unit"`
	SupplierID    string       `json:"supplier_id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type InspectionItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Standard    string  `json:"standard"`
	Unit        string  `json:"unit"`
	MinValue    float64 `json:"min_value"`
	MaxValue    float64 `json:"max_value"`
	IsNumeric   bool    `json:"is_numeric"`
	SPCEnabled  bool    `json:"spc_enabled"`
	Description string  `json:"description"`
}

type Standard struct {
	ID               string           `json:"id"`
	MaterialID       string           `json:"material_id"`
	Type             InspectionType   `json:"type"`
	InspectionLevels []InspectionItem `json:"inspection_items"`
	AQL              float64          `json:"aql"`
	InspectionLevel  string           `json:"inspection_level"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}
