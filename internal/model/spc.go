package model

import "time"

type SPCControlChart struct {
	ID          string          `json:"id"`
	MaterialID  string          `json:"material_id"`
	ItemID      string          `json:"item_id"`
	ItemName    string          `json:"item_name"`
	ChartType   string          `json:"chart_type"`
	XBarCL      float64         `json:"xbar_cl"`
	XBarUCL     float64         `json:"xbar_ucl"`
	XBarLCL     float64         `json:"xbar_lcl"`
	RCL         float64         `json:"r_cl"`
	RUCL        float64         `json:"r_ucl"`
	RLCL        float64         `json:"r_lcl"`
	CL          float64         `json:"cl"`
	UCL         float64         `json:"ucl"`
	LCL         float64         `json:"lcl"`
	Subgroups   []SubgroupData  `json:"subgroups"`
	LastUpdated time.Time       `json:"last_updated"`
}

type SubgroupData struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	XBar      float64   `json:"x_bar"`
	R         float64   `json:"r"`
	SampleSize int      `json:"sample_size"`
	Values    []float64 `json:"values"`
	IsOutOfControl bool `json:"is_out_of_control"`
}
