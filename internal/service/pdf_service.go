package service

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GenerateInspectionReport(recordID string) ([]byte, error) {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return nil, ErrRecordNotFound
	}

	lot, _ := store.GlobalStore.GetLot(rec.LotID)
	material, _ := store.GlobalStore.GetMaterial(rec.MaterialID)
	supplier, _ := store.GlobalStore.GetSupplier(rec.SupplierID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Inspection Report")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Basic Info")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)

	basicInfo := [][]string{
		{"Record ID", rec.ID},
		{"Lot No", rec.LotNo},
		{"Type", string(rec.Type)},
		{"Status", string(rec.Status)},
		{"Version", fmt.Sprintf("v%d", rec.Version)},
		{"Inspector", rec.InspectorName},
		{"Created At", rec.CreatedAt.Format("2006-01-02 15:04:05")},
	}

	if material != nil {
		basicInfo = append(basicInfo, []string{"Material Code", material.Code})
		basicInfo = append(basicInfo, []string{"Material Name", material.Name})
	}
	if lot != nil {
		basicInfo = append(basicInfo, []string{"Quantity", fmt.Sprintf("%d", lot.Quantity)})
		basicInfo = append(basicInfo, []string{"Production Date", lot.ProductionDate.Format("2006-01-02")})
		basicInfo = append(basicInfo, []string{"Shift", lot.Shift})
		basicInfo = append(basicInfo, []string{"Production Line", lot.ProductionLine})
	}
	if supplier != nil {
		basicInfo = append(basicInfo, []string{"Supplier", supplier.Name})
	}

	for _, info := range basicInfo {
		pdf.Cell(40, 6, info[0]+":")
		pdf.Cell(60, 6, info[1])
		pdf.Ln(6)
	}

	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Inspection Items")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 7, "Item Name")
	pdf.Cell(50, 7, "Standard")
	pdf.Cell(30, 7, "Actual Value")
	pdf.Cell(25, 7, "Judgment")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)

	for _, item := range rec.Items {
		pdf.Cell(40, 6, item.ItemName)
		pdf.Cell(50, 6, item.Standard)
		pdf.Cell(30, 6, item.ActualValue)
		if item.Judgment == model.JudgmentPass {
			pdf.SetTextColor(0, 128, 0)
		} else {
			pdf.SetTextColor(255, 0, 0)
		}
		pdf.Cell(25, 6, string(item.Judgment))
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(6)
	}

	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Summary")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)

	summaryInfo := [][]string{
		{"Total Sample Size", fmt.Sprintf("%d", rec.TotalSampleSize)},
		{"Defect Count", fmt.Sprintf("%d", rec.DefectCount)},
		{"Final Judgment", string(rec.FinalJudgment)},
		{"Disposition", string(rec.Disposition)},
	}

	for _, info := range summaryInfo {
		pdf.Cell(40, 6, info[0]+":")
		if info[0] == "Final Judgment" {
			if info[1] == string(model.JudgmentPass) {
				pdf.SetTextColor(0, 128, 0)
			} else {
				pdf.SetTextColor(255, 0, 0)
			}
		}
		pdf.Cell(60, 6, info[1])
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(6)
	}

	if len(rec.Remarks) > 0 {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 8, "Remarks")
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 10)

		for i, remark := range rec.Remarks {
			pdf.Cell(10, 6, fmt.Sprintf("%d.", i+1))
			pdf.Cell(120, 6, remark.Content)
			pdf.Cell(30, 6, remark.Operator)
			pdf.Ln(6)
		}
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *PDFService) GetReportFileName(recordID string, lotNo string) string {
	return fmt.Sprintf("inspection_report_%s_%s.pdf", lotNo, recordID[:8])
}
