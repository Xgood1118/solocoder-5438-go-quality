package model

import "time"

type Nonconformance struct {
	ID             string    `json:"id"`
	NCRNo          string    `json:"ncr_no"`
	RecordID       string    `json:"record_id"`
	LotID          string    `json:"lot_id"`
	MaterialID     string    `json:"material_id"`
	SupplierID     string    `json:"supplier_id"`
	Description    string    `json:"description"`
	Severity       string    `json:"severity"`
	Quantity       int       `json:"quantity"`
	Disposition    DispositionType `json:"disposition"`
	RootCause      string    `json:"root_cause"`
	CorrectiveAction string  `json:"corrective_action"`
	PreventiveAction string  `json:"preventive_action"`
	Responsible    string    `json:"responsible"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `json:"status"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
}

type ReturnOrder struct {
	ID           string    `json:"id"`
	ReturnNo     string    `json:"return_no"`
	SupplierID   string    `json:"supplier_id"`
	RecordID     string    `json:"record_id"`
	LotID        string    `json:"lot_id"`
	MaterialID   string    `json:"material_id"`
	Quantity     int       `json:"quantity"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}
