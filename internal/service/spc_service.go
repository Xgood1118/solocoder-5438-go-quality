package service

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type SPCService struct{}

func NewSPCService() *SPCService {
	return &SPCService{}
}

func (s *SPCService) CalculateControlLimits(values []float64) (cl, ucl, lcl float64) {
	n := len(values)
	if n < 2 {
		return 0, 0, 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)

	variance := 0.0
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(variance / float64(n-1))

	cl = mean
	ucl = mean + 3*stdDev
	lcl = mean - 3*stdDev

	return cl, ucl, lcl
}

func (s *SPCService) CalculateXBarR(subgroups [][]float64) (xBarCL, xBarUCL, xBarLCL, rCL, rUCL, rLCL float64) {
	if len(subgroups) < 2 {
		return 0, 0, 0, 0, 0, 0
	}

	subgroupSize := len(subgroups[0])

	xBars := make([]float64, len(subgroups))
	ranges := make([]float64, len(subgroups))

	for i, sg := range subgroups {
		sum := 0.0
		minVal := math.Inf(1)
		maxVal := math.Inf(-1)
		for _, v := range sg {
			sum += v
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
		xBars[i] = sum / float64(len(sg))
		ranges[i] = maxVal - minVal
	}

	xBarSum := 0.0
	for _, xb := range xBars {
		xBarSum += xb
	}
	xBarCL = xBarSum / float64(len(xBars))

	rSum := 0.0
	for _, r := range ranges {
		rSum += r
	}
	rCL = rSum / float64(len(ranges))

	A2 := getA2(subgroupSize)
	D3 := getD3(subgroupSize)
	D4 := getD4(subgroupSize)

	xBarUCL = xBarCL + A2*rCL
	xBarLCL = xBarCL - A2*rCL
	rUCL = D4 * rCL
	rLCL = D3 * rCL

	return xBarCL, xBarUCL, xBarLCL, rCL, rUCL, rLCL
}

func getA2(n int) float64 {
	a2Table := map[int]float64{
		2:  1.880,
		3:  1.023,
		4:  0.729,
		5:  0.577,
		6:  0.483,
		7:  0.419,
		8:  0.373,
		9:  0.337,
		10: 0.308,
	}
	if n <= 1 {
		return 0
	}
	if n > 10 {
		return 3 / math.Sqrt(float64(n))
	}
	return a2Table[n]
}

func getD3(n int) float64 {
	d3Table := map[int]float64{
		2: 0.000,
		3: 0.000,
		4: 0.000,
		5: 0.000,
		6: 0.000,
		7: 0.076,
		8: 0.136,
		9: 0.184,
		10: 0.223,
	}
	if n < 2 {
		return 0
	}
	if n > 10 {
		return 0
	}
	return d3Table[n]
}

func getD4(n int) float64 {
	d4Table := map[int]float64{
		2: 3.267,
		3: 2.574,
		4: 2.282,
		5: 2.114,
		6: 2.004,
		7: 1.924,
		8: 1.864,
		9: 1.816,
		10: 1.777,
	}
	if n < 2 {
		return 0
	}
	if n > 10 {
		return 3/math.Sqrt(float64(n)) + 1
	}
	return d4Table[n]
}

func (s *SPCService) IsOutOfControl(value float64, ucl, lcl float64) bool {
	return value > ucl || value < lcl
}

func (s *SPCService) UpdateSPCChart(materialID string, itemID string, itemName string) (*model.SPCControlChart, error) {
	records := store.GlobalStore.ListRecords()

	var allValues []float64
	var subgroupValues [][]float64
	currentSubgroup := make([]float64, 0, 5)

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.Before(records[j].CreatedAt)
	})

	for _, rec := range records {
		if rec.MaterialID != materialID || rec.Status != model.StatusApproved {
			continue
		}
		for _, item := range rec.Items {
			if item.ItemID == itemID && item.IsNumeric {
				allValues = append(allValues, item.NumericValue)
				currentSubgroup = append(currentSubgroup, item.NumericValue)
				if len(currentSubgroup) >= 5 {
					sg := make([]float64, len(currentSubgroup))
					copy(sg, currentSubgroup)
					subgroupValues = append(subgroupValues, sg)
					currentSubgroup = make([]float64, 0, 5)
				}
			}
		}
	}

	if len(allValues) < 25 {
		if len(allValues) < 2 {
			return nil, nil
		}
		cl, ucl, lcl := s.CalculateControlLimits(allValues)

		chart := &model.SPCControlChart{
			ID:          uuid.New().String(),
			MaterialID:  materialID,
			ItemID:      itemID,
			ItemName:    itemName,
			ChartType:   "xbar_r",
			CL:          cl,
			UCL:         ucl,
			LCL:         lcl,
			XBarCL:      cl,
			XBarUCL:     ucl,
			XBarLCL:     lcl,
			LastUpdated: time.Now(),
		}

		for _, v := range allValues {
			sg := model.SubgroupData{
				ID:            uuid.New().String(),
				Timestamp:     time.Now(),
				XBar:          v,
				R:             0,
				SampleSize:    1,
				Values:        []float64{v},
				IsOutOfControl: s.IsOutOfControl(v, ucl, lcl),
			}
			chart.Subgroups = append(chart.Subgroups, sg)
		}

		existing, ok := store.GlobalStore.GetSPCChart(materialID, itemID)
		if ok {
			chart.ID = existing.ID
		}
		store.GlobalStore.SaveSPCChart(chart)

		return chart, nil
	}

	xBarCL, xBarUCL, xBarLCL, rCL, rUCL, rLCL := s.CalculateXBarR(subgroupValues)

	chart := &model.SPCControlChart{
		ID:          uuid.New().String(),
		MaterialID:  materialID,
		ItemID:      itemID,
		ItemName:    itemName,
		ChartType:   "xbar_r",
		CL:          xBarCL,
		UCL:         xBarUCL,
		LCL:         xBarLCL,
		XBarCL:      xBarCL,
		XBarUCL:     xBarUCL,
		XBarLCL:     xBarLCL,
		RCL:         rCL,
		RUCL:        rUCL,
		RLCL:        rLCL,
		LastUpdated: time.Now(),
	}

	for i, sg := range subgroupValues {
		sum := 0.0
		minVal := math.Inf(1)
		maxVal := math.Inf(-1)
		for _, v := range sg {
			sum += v
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
		xBar := sum / float64(len(sg))
		r := maxVal - minVal

		subgroup := model.SubgroupData{
			ID:            uuid.New().String(),
			Timestamp:     time.Now().Add(time.Duration(-i) * time.Hour),
			XBar:          xBar,
			R:             r,
			SampleSize:    len(sg),
			Values:        sg,
			IsOutOfControl: s.IsOutOfControl(xBar, xBarUCL, xBarLCL) || s.IsOutOfControl(r, rUCL, 0),
		}
		chart.Subgroups = append(chart.Subgroups, subgroup)
	}

	existing, ok := store.GlobalStore.GetSPCChart(materialID, itemID)
	if ok {
		chart.ID = existing.ID
	}
	store.GlobalStore.SaveSPCChart(chart)

	return chart, nil
}

func (s *SPCService) GetSPCChart(materialID, itemID string) (*model.SPCControlChart, bool) {
	return store.GlobalStore.GetSPCChart(materialID, itemID)
}

func (s *SPCService) CheckAlarms(chart *model.SPCControlChart) []model.SubgroupData {
	var alarms []model.SubgroupData
	for _, sg := range chart.Subgroups {
		if sg.IsOutOfControl {
			alarms = append(alarms, sg)
		}
	}
	return alarms
}
