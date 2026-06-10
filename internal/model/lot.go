package model

import "time"

type Lot struct {
	ID             string    `json:"id"`
	LotNo          string    `json:"lot_no"`
	MaterialID     string    `json:"material_id"`
	Quantity       int       `json:"quantity"`
	ProductionDate time.Time `json:"production_date"`
	Shift          string    `json:"shift"`
	Worker         string    `json:"worker"`
	ProductionLine string    `json:"production_line"`
	Equipment      string    `json:"equipment"`
	RawMaterialLot string    `json:"raw_material_lot"`
	SupplierID     string    `json:"supplier_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
