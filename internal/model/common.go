package model

import "time"

type InspectionType string

const (
	TypeIQC  InspectionType = "IQC"
	TypeIPQC InspectionType = "IPQC"
	TypeOQC  InspectionType = "OQC"
	TypeOBC  InspectionType = "OBC"
)

type MaterialType string

const (
	MaterialMetal     MaterialType = "metal"
	MaterialPlastic   MaterialType = "plastic"
	MaterialElectronic MaterialType = "electronic"
	MaterialChemical  MaterialType = "chemical"
	MaterialPaper     MaterialType = "paper"
)

type JudgmentResult string

const (
	JudgmentPass    JudgmentResult = "Pass"
	JudgmentFail    JudgmentResult = "Fail"
	JudgmentReject  JudgmentResult = "Reject"
)

type DispositionType string

const (
	DispositionConcession DispositionType = "concession"
	DispositionRework     DispositionType = "rework"
	DispositionScrap      DispositionType = "scrap"
	DispositionReturn     DispositionType = "return"
)

type RecordStatus string

const (
	StatusDraft     RecordStatus = "Draft"
	StatusSubmitted RecordStatus = "Submitted"
	StatusApproved  RecordStatus = "Approved"
	StatusRejected  RecordStatus = "Rejected"
	StatusClosed    RecordStatus = "Closed"
)

type IPQCType string

const (
	IPQCFirst  IPQCType = "first_article"
	IPQCPatrol IPQCType = "patrol"
	IPQCLast   IPQCType = "last_article"
)

type ApprovalStatus string

const (
	ApprovalPending    ApprovalStatus = "pending"
	ApprovalQA         ApprovalStatus = "qa_approved"
	ApprovalDirector   ApprovalStatus = "director_approved"
	ApprovalRejected   ApprovalStatus = "rejected"
)

type AuditLog struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Operator  string    `json:"operator"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *Pagination) Default() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
