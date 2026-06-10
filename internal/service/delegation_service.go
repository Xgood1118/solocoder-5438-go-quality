package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

type DelegationService struct{}

func NewDelegationService() *DelegationService {
	return &DelegationService{}
}

func (s *DelegationService) CreateDelegation(delegatorID string, delegateeID string, role string, days int) (*model.Delegation, error) {
	if days <= 0 || days > 30 {
		return nil, errors.New("invalid delegation duration, max 30 days")
	}

	s.expireOldDelegations(delegatorID)

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, days)

	delegation := &model.Delegation{
		ID:            uuid.New().String(),
		DelegatorID:   delegatorID,
		DelegateeID:   delegateeID,
		DelegatorRole: role,
		StartDate:     startDate,
		EndDate:       endDate,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}

	store.GlobalStore.SaveDelegation(delegation)

	transferService := NewRecordService()
	_, err := transferService.TransferRecords(delegatorID, delegateeID, "system")
	if err != nil {
		return nil, err
	}

	return delegation, nil
}

func (s *DelegationService) GetActiveDelegation(delegatorID string) (*model.Delegation, bool) {
	return store.GlobalStore.GetActiveDelegation(delegatorID)
}

func (s *DelegationService) RevokeDelegation(delegatorID string) error {
	d, ok := store.GlobalStore.GetActiveDelegation(delegatorID)
	if !ok {
		return errors.New("no active delegation")
	}
	d.IsActive = false
	return nil
}

func (s *DelegationService) expireOldDelegations(delegatorID string) {
}

func (s *DelegationService) GetActualApprover(approverID string, role string) string {
	delegation, ok := store.GlobalStore.GetActiveDelegation(approverID)
	if ok && delegation.DelegatorRole == role {
		return delegation.DelegateeID
	}
	return approverID
}
