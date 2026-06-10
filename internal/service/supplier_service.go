package service

import (
	"time"

	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type SupplierScoreService struct{}

func NewSupplierScoreService() *SupplierScoreService {
	return &SupplierScoreService{}
}

func (s *SupplierScoreService) CalculateMonthlyScore(supplierID string, yearMonth string) (*model.SupplierScore, error) {
	supplier, ok := store.GlobalStore.GetSupplier(supplierID)
	if !ok {
		return nil, ErrRecordNotFound
	}

	records := store.GlobalStore.ListRecords()

	var iqcRecords []*model.InspectionRecord
	var ncRecords []*model.InspectionRecord
	var returnCount int

	for _, rec := range records {
		if rec.SupplierID != supplierID || rec.Type != model.TypeIQC {
			continue
		}

		recYearMonth := rec.CreatedAt.Format("200601")
		if recYearMonth != yearMonth {
			continue
		}

		iqcRecords = append(iqcRecords, rec)

		if rec.FinalJudgment == model.JudgmentFail || rec.FinalJudgment == model.JudgmentReject {
			ncRecords = append(ncRecords, rec)
		}

		if rec.IsReturnGoods {
			returnCount++
		}
	}

	totalCount := len(iqcRecords)
	if totalCount == 0 {
		return &model.SupplierScore{
			ID:                 uuid.New().String(),
			SupplierID:         supplierID,
			YearMonth:          yearMonth,
			IQCPassRate:        100,
			ReturnRate:         0,
			ResponseSpeedScore: 80,
			TotalScore:         85,
			Level:              "normal",
			CreatedAt:          time.Now(),
		}, nil
	}

	passCount := 0
	for _, rec := range iqcRecords {
		if rec.FinalJudgment == model.JudgmentPass {
			passCount++
		}
	}

	passRate := float64(passCount) / float64(totalCount) * 100
	returnRate := float64(returnCount) / float64(totalCount) * 100

	passRateScore := passRate * 0.5

	returnRateScore := 100.0
	if returnRate > 0 {
		returnRateScore = 100 - returnRate*2
		if returnRateScore < 0 {
			returnRateScore = 0
		}
	}
	returnRateScore *= 0.3

	responseScore := 80.0 * 0.2

	totalScore := passRateScore + returnRateScore + responseScore

	level := "normal"
	if totalScore < 40 {
		level = "red"
	} else if totalScore < 60 {
		level = "yellow"
	}

	score := &model.SupplierScore{
		ID:                 uuid.New().String(),
		SupplierID:         supplierID,
		YearMonth:          yearMonth,
		IQCPassRate:        passRate,
		ReturnRate:         returnRate,
		ResponseSpeedScore: responseScore * 5,
		TotalScore:         totalScore,
		Level:              level,
		CreatedAt:          time.Now(),
	}

	store.GlobalStore.SaveSupplierScore(score)

	supplier.Rating = totalScore
	if level == "red" {
		supplier.Status = "red"
	} else if level == "yellow" {
		supplier.Status = "yellow"
	} else {
		supplier.Status = "normal"
	}
	supplier.UpdatedAt = time.Now()

	return score, nil
}

func (s *SupplierScoreService) GetSupplierScores(supplierID string) []*model.SupplierScore {
	return store.GlobalStore.GetSupplierScores(supplierID)
}

func (s *SupplierScoreService) CalculateAllSuppliers(yearMonth string) []*model.SupplierScore {
	suppliers := store.GlobalStore.ListSuppliers()
	var results []*model.SupplierScore

	for _, sup := range suppliers {
		score, err := s.CalculateMonthlyScore(sup.ID, yearMonth)
		if err == nil {
			results = append(results, score)
		}
	}

	return results
}

func (s *SupplierScoreService) GetFocusSuppliers() []*model.Supplier {
	suppliers := store.GlobalStore.ListSuppliers()
	var result []*model.Supplier

	for _, sup := range suppliers {
		if sup.Status == "red" || sup.Status == "yellow" {
			result = append(result, sup)
		}
	}

	return result
}
