package aql

import (
	"errors"
	"sort"
)

type SampleResult struct {
	SampleSize int
	AcceptNum  int
	RejectNum  int
}

var (
	batchRanges = []int{25, 50, 90, 150, 280, 500, 1200, 3200, 10000, 35000}

	aqlValues = []float64{0.065, 0.15, 0.25, 0.4, 0.65, 1.0, 1.5, 2.5, 4.0, 6.5}

	sampleSizeTableLevelII = []int{5, 8, 13, 20, 32, 50, 80, 125, 200, 315}

	acceptTableLevelII = [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 2},
		{0, 0, 0, 0, 0, 0, 1, 2, 3, 5},
		{0, 0, 0, 0, 0, 1, 2, 3, 5, 8},
		{0, 0, 0, 0, 1, 2, 3, 5, 8, 12},
		{0, 0, 0, 1, 2, 3, 5, 8, 12, 18},
		{0, 0, 1, 2, 3, 5, 8, 12, 18, 26},
		{0, 1, 2, 3, 5, 8, 12, 18, 26, 40},
		{1, 2, 3, 5, 8, 12, 18, 26, 40, 63},
		{2, 3, 5, 8, 12, 18, 26, 40, 63, 99},
		{3, 5, 8, 12, 18, 26, 40, 63, 99, 150},
	}
)

func getBatchIndex(batchSize int) int {
	if batchSize <= 0 {
		return 0
	}
	idx := sort.SearchInts(batchRanges, batchSize)
	if idx >= len(batchRanges) {
		idx = len(batchRanges) - 1
	}
	return idx
}

func getAQLIndex(aql float64) (int, error) {
	for i, v := range aqlValues {
		if v == aql {
			return i, nil
		}
	}
	return -1, errors.New("invalid AQL value")
}

func CalculateSample(batchSize int, aql float64, inspectionLevel string) (*SampleResult, error) {
	batchIdx := getBatchIndex(batchSize)
	aqlIdx, err := getAQLIndex(aql)
	if err != nil {
		return nil, err
	}

	var sampleSize int
	var acceptNum int

	switch inspectionLevel {
	case "I":
		sampleSize = sampleSizeTableLevelII[batchIdx] / 2
		if sampleSize < 2 {
			sampleSize = 2
		}
		acceptNum = acceptTableLevelII[batchIdx][aqlIdx] / 2
	case "II", "":
		sampleSize = sampleSizeTableLevelII[batchIdx]
		acceptNum = acceptTableLevelII[batchIdx][aqlIdx]
	case "III":
		sampleSize = sampleSizeTableLevelII[batchIdx] * 2
		acceptNum = acceptTableLevelII[batchIdx][aqlIdx] * 2
	default:
		return nil, errors.New("invalid inspection level")
	}

	if sampleSize > batchSize {
		sampleSize = batchSize
	}

	return &SampleResult{
		SampleSize: sampleSize,
		AcceptNum:  acceptNum,
		RejectNum:  acceptNum + 1,
	}, nil
}

func Judge(sampleSize int, defectCount int, acceptNum int) bool {
	return defectCount <= acceptNum
}

func GetBatchRanges() []int {
	result := make([]int, len(batchRanges))
	copy(result, batchRanges)
	return result
}

func GetAQLValues() []float64 {
	result := make([]float64, len(aqlValues))
	copy(result, aqlValues)
	return result
}

func RoundBatchSize(batchSize int) int {
	if batchSize <= 0 {
		return 0
	}
	idx := sort.SearchInts(batchRanges, batchSize)
	if idx >= len(batchRanges) {
		return batchRanges[len(batchRanges)-1]
	}
	return batchRanges[idx]
}
