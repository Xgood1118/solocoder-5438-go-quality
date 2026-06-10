package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type OCRService struct{}

func NewOCRService() *OCRService {
	return &OCRService{}
}

func (s *OCRService) RecognizeMark(imageData []byte) (string, error) {
	if len(imageData) == 0 {
		return "", errors.New("empty image data")
	}

	if rand.Float32() < 0.8 {
		markText := fmt.Sprintf("LOT-%s-%04d", time.Now().Format("20060102"), rand.Intn(9999))
		return markText, nil
	}

	return "", errors.New("OCR recognition failed")
}

func (s *OCRService) RecognizeWithFallback(imageData []byte, manualInput string) (string, error) {
	result, err := s.RecognizeMark(imageData)
	if err == nil {
		return result, nil
	}

	if manualInput != "" {
		return manualInput, nil
	}

	return "", errors.New("OCR failed and no manual input provided")
}
