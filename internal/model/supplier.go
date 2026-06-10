package model

import "time"

type Supplier struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Rating    float64   `json:"rating"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SupplierScore struct {
	ID                 string    `json:"id"`
	SupplierID         string    `json:"supplier_id"`
	YearMonth          string    `json:"year_month"`
	IQCPassRate        float64   `json:"iqc_pass_rate"`
	ReturnRate         float64   `json:"return_rate"`
	ResponseSpeedScore float64   `json:"response_speed_score"`
	TotalScore         float64   `json:"total_score"`
	Level              string    `json:"level"`
	CreatedAt          time.Time `json:"created_at"`
}
