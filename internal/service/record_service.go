package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/store"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrRecordNotFound          = errors.New("record not found")
	ErrNoPermission            = errors.New("no permission")
)

type RecordService struct{}

func NewRecordService() *RecordService {
	return &RecordService{}
}

func (s *RecordService) canEdit(status model.RecordStatus, role string) bool {
	if role == "admin" || role == "qa_manager" || role == "director" {
		return true
	}
	return status == model.StatusDraft
}

func (s *RecordService) canAddRemark(status model.RecordStatus) bool {
	return status != model.StatusClosed
}

func (s *RecordService) transitionStatus(rec *model.InspectionRecord, targetStatus model.RecordStatus, operator string, remark string) error {
	validTransitions := map[model.RecordStatus][]model.RecordStatus{
		model.StatusDraft:     {model.StatusSubmitted, model.StatusClosed},
		model.StatusSubmitted: {model.StatusApproved, model.StatusRejected, model.StatusDraft},
		model.StatusApproved:  {model.StatusClosed},
		model.StatusRejected:  {model.StatusClosed},
	}

	validTargets, ok := validTransitions[rec.Status]
	if !ok {
		return ErrInvalidStatusTransition
	}

	valid := false
	for _, t := range validTargets {
		if t == targetStatus {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidStatusTransition
	}

	oldStatus := rec.Status
	rec.Status = targetStatus

	now := time.Now()
	switch targetStatus {
	case model.StatusSubmitted:
		rec.SubmittedAt = &now
	case model.StatusApproved:
		rec.ApprovedAt = &now
	case model.StatusClosed:
		rec.ClosedAt = &now
	}

	store.GlobalStore.UpdateRecord(rec)

	auditLog := &model.AuditLog{
		ID:        uuid.New().String(),
		RecordID:  rec.ID,
		Action:    "status_change",
		OldStatus: string(oldStatus),
		NewStatus: string(targetStatus),
		Operator:  operator,
		Remark:    remark,
		CreatedAt: now,
	}
	store.GlobalStore.SaveAuditLog(auditLog)

	return nil
}

func (s *RecordService) SubmitRecord(recordID string, operator string) error {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return ErrRecordNotFound
	}
	return s.transitionStatus(rec, model.StatusSubmitted, operator, "提交检验记录")
}

func (s *RecordService) ApproveRecord(recordID string, operator string, role string, remark string) error {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return ErrRecordNotFound
	}

	if rec.Disposition != model.DispositionConcession {
		return s.transitionStatus(rec, model.StatusApproved, operator, remark)
	}

	approval, ok := store.GlobalStore.GetApproval(rec.ApprovalID)
	if !ok {
		return errors.New("approval not found")
	}

	now := time.Now()
	switch role {
	case "qa_manager":
		if approval.Status != model.ApprovalPending {
			return errors.New("invalid approval status for QA approval")
		}
		approval.Status = model.ApprovalQA
		approval.QAApproverID = operator
		approval.QAApproverName = operator
		approval.QARemark = remark
		approval.QAApprovedAt = &now
		store.GlobalStore.UpdateApproval(approval)
	case "director":
		if approval.Status != model.ApprovalQA {
			return errors.New("invalid approval status for director approval")
		}
		approval.Status = model.ApprovalDirector
		approval.DirectorApproverID = operator
		approval.DirectorApproverName = operator
		approval.DirectorRemark = remark
		approval.DirectorApprovedAt = &now
		store.GlobalStore.UpdateApproval(approval)

		return s.transitionStatus(rec, model.StatusApproved, operator, "厂长终审通过")
	default:
		return ErrNoPermission
	}

	return nil
}

func (s *RecordService) RejectRecord(recordID string, operator string, role string, remark string) error {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return ErrRecordNotFound
	}

	if rec.Disposition == model.DispositionConcession {
		approval, ok := store.GlobalStore.GetApproval(rec.ApprovalID)
		if ok {
			approval.Status = model.ApprovalRejected
			store.GlobalStore.UpdateApproval(approval)
		}
	}

	return s.transitionStatus(rec, model.StatusRejected, operator, remark)
}

func (s *RecordService) CloseRecord(recordID string, operator string, remark string) error {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return ErrRecordNotFound
	}

	if rec.Status == model.StatusRejected {
		nc := &model.Nonconformance{
			ID:          uuid.New().String(),
			NCRNo:       generateNCRNo(),
			RecordID:    rec.ID,
			LotID:       rec.LotID,
			MaterialID:  rec.MaterialID,
			SupplierID:  rec.SupplierID,
			Description: remark,
			Severity:    "major",
			Quantity:    rec.TotalSampleSize,
			Disposition: rec.Disposition,
			Status:      "open",
			CreatedBy:   operator,
			CreatedAt:   time.Now(),
		}
		store.GlobalStore.SaveNonconformance(nc)

		if rec.Disposition == model.DispositionReturn {
			ro := &model.ReturnOrder{
				ID:          uuid.New().String(),
				ReturnNo:    generateReturnNo(),
				SupplierID:  rec.SupplierID,
				RecordID:    rec.ID,
				LotID:       rec.LotID,
				MaterialID:  rec.MaterialID,
				Quantity:    rec.TotalSampleSize,
				Reason:      remark,
				Status:      "created",
				CreatedBy:   operator,
				CreatedAt:   time.Now(),
			}
			store.GlobalStore.SaveReturnOrder(ro)
			rec.IsReturnGoods = true
		}
	}

	return s.transitionStatus(rec, model.StatusClosed, operator, remark)
}

func (s *RecordService) AddRemark(recordID string, content string, operator string) error {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return ErrRecordNotFound
	}

	if !s.canAddRemark(rec.Status) {
		return errors.New("cannot add remark to closed record")
	}

	remark := model.Remark{
		ID:        uuid.New().String(),
		Content:   content,
		Operator:  operator,
		CreatedAt: time.Now(),
	}
	rec.Remarks = append(rec.Remarks, remark)
	store.GlobalStore.UpdateRecord(rec)

	auditLog := &model.AuditLog{
		ID:        uuid.New().String(),
		RecordID:  rec.ID,
		Action:    "add_remark",
		OldStatus: string(rec.Status),
		NewStatus: string(rec.Status),
		Operator:  operator,
		Remark:    content,
		CreatedAt: time.Now(),
	}
	store.GlobalStore.SaveAuditLog(auditLog)

	return nil
}

func (s *RecordService) CreateReinspection(originalRecordID string, operator string) (*model.InspectionRecord, error) {
	orig, ok := store.GlobalStore.GetRecord(originalRecordID)
	if !ok {
		return nil, ErrRecordNotFound
	}

	newVersion := orig.Version + 1
	newRecord := &model.InspectionRecord{
		ID:               uuid.New().String(),
		LotID:            orig.LotID,
		LotNo:            orig.LotNo,
		MaterialID:       orig.MaterialID,
		SupplierID:       orig.SupplierID,
		Type:             orig.Type,
		IPQCType:         orig.IPQCType,
		InspectorID:      operator,
		InspectorName:    operator,
		Status:           model.StatusDraft,
		TotalSampleSize:  orig.TotalSampleSize,
		DefectCount:      0,
		FinalJudgment:    "",
		Disposition:      "",
		Items:            make([]model.InspectionItemResult, len(orig.Items)),
		Version:          newVersion,
		ParentRecordID:   orig.ID,
		Remarks:          []model.Remark{},
		IsFullInspection: orig.IsFullInspection,
		CreatedAt:        time.Now(),
	}

	for i, item := range orig.Items {
		newRecord.Items[i] = model.InspectionItemResult{
			ItemID:     item.ItemID,
			ItemName:   item.ItemName,
			Standard:   item.Standard,
			IsNumeric:  item.IsNumeric,
			Judgment:   "",
		}
	}

	store.GlobalStore.SaveRecord(newRecord)

	return newRecord, nil
}

func (s *RecordService) TransferRecords(fromInspectorID string, toInspectorID string, operator string) (int, error) {
	records := store.GlobalStore.GetRecordsByInspector(fromInspectorID)
	count := 0

	for _, rec := range records {
		rec.InspectorID = toInspectorID

		remark := model.Remark{
			ID:        uuid.New().String(),
			Content:   "检验单从 " + fromInspectorID + " 转移到 " + toInspectorID,
			Operator:  operator,
			CreatedAt: time.Now(),
		}
		rec.Remarks = append(rec.Remarks, remark)

		store.GlobalStore.UpdateRecord(rec)

		auditLog := &model.AuditLog{
			ID:        uuid.New().String(),
			RecordID:  rec.ID,
			Action:    "transfer",
			OldStatus: fromInspectorID,
			NewStatus: toInspectorID,
			Operator:  operator,
			Remark:    "检验员转移",
			CreatedAt: time.Now(),
		}
		store.GlobalStore.SaveAuditLog(auditLog)

		count++
	}

	return count, nil
}

func (s *RecordService) CreateConcessionApproval(recordID string, operator string) (*model.Approval, error) {
	rec, ok := store.GlobalStore.GetRecord(recordID)
	if !ok {
		return nil, ErrRecordNotFound
	}

	approval := &model.Approval{
		ID:        uuid.New().String(),
		RecordID:  recordID,
		Type:      "concession",
		Status:    model.ApprovalPending,
		CreatedAt: time.Now(),
	}

	store.GlobalStore.SaveApproval(approval)
	rec.ApprovalID = approval.ID
	rec.Disposition = model.DispositionConcession
	store.GlobalStore.UpdateRecord(rec)

	return approval, nil
}

func generateNCRNo() string {
	now := time.Now()
	return "NCR-" + now.Format("200601") + "-" + uuid.New().String()[:6]
}

func generateReturnNo() string {
	now := time.Now()
	return "RTN-" + now.Format("200601") + "-" + uuid.New().String()[:6]
}
