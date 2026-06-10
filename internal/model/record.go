package model

import "time"

type InspectionRecord struct {
	ID               string                `json:"id"`
	LotID            string                `json:"lot_id"`
	LotNo            string                `json:"lot_no"`
	MaterialID       string                `json:"material_id"`
	SupplierID       string                `json:"supplier_id"`
	Type             InspectionType        `json:"type"`
	IPQCType         IPQCType              `json:"ipqc_type,omitempty"`
	InspectorID      string                `json:"inspector_id"`
	InspectorName    string                `json:"inspector_name"`
	Status           RecordStatus          `json:"status"`
	TotalSampleSize  int                   `json:"total_sample_size"`
	DefectCount      int                   `json:"defect_count"`
	FinalJudgment    JudgmentResult        `json:"final_judgment"`
	Disposition      DispositionType       `json:"disposition"`
	Items            []InspectionItemResult `json:"items"`
	Version          int                   `json:"version"`
	ParentRecordID   string                `json:"parent_record_id,omitempty"`
	ApprovalID       string                `json:"approval_id,omitempty"`
	Remarks          []Remark              `json:"remarks"`
	IsFullInspection bool                  `json:"is_full_inspection"`
	IsReturnGoods    bool                  `json:"return_goods"`
	OCRMarkText      string                `json:"ocr_mark_text,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	SubmittedAt      *time.Time            `json:"submitted_at,omitempty"`
	ApprovedAt       *time.Time            `json:"approved_at,omitempty"`
	ClosedAt         *time.Time            `json:"closed_at,omitempty"`
}

type InspectionItemResult struct {
	ItemID      string         `json:"item_id"`
	ItemName    string         `json:"item_name"`
	Standard    string         `json:"standard"`
	ActualValue string         `json:"actual_value"`
	NumericValue float64       `json:"numeric_value"`
	IsNumeric   bool           `json:"is_numeric"`
	Judgment    JudgmentResult `json:"judgment"`
	DefectDesc  string         `json:"defect_desc"`
}

type Remark struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Operator  string    `json:"operator"`
	CreatedAt time.Time `json:"created_at"`
}

type Approval struct {
	ID              string         `json:"id"`
	RecordID        string         `json:"record_id"`
	Type            string         `json:"type"`
	Status          ApprovalStatus `json:"status"`
	QAApproverID    string         `json:"qa_approver_id,omitempty"`
	QAApproverName  string         `json:"qa_approver_name,omitempty"`
	QARemark        string         `json:"qa_remark,omitempty"`
	QAApprovedAt    *time.Time     `json:"qa_approved_at,omitempty"`
	DirectorApproverID   string    `json:"director_approver_id,omitempty"`
	DirectorApproverName string    `json:"director_approver_name,omitempty"`
	DirectorRemark  string         `json:"director_remark,omitempty"`
	DirectorApprovedAt *time.Time  `json:"director_approved_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}
