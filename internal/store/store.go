package store

import (
	"sync"
	"sync/atomic"
	"time"

	"qc-system/internal/model"
)

func timeNow() time.Time {
	return time.Now()
}

type Store struct {
	Materials        sync.Map
	Suppliers        sync.Map
	Inspectors       sync.Map
	Lots             sync.Map
	Standards        sync.Map
	Records          sync.Map
	Nonconformances  sync.Map
	ReturnOrders     sync.Map
	AuditLogs        sync.Map
	Approvals        sync.Map
	Delegations      sync.Map
	SupplierScores   sync.Map
	SPCCharts        sync.Map

	materialCount int64
	supplierCount int64
	inspectorCount int64
	lotCount      int64
	standardCount int64
	recordCount   int64
	ncCount       int64
	returnCount   int64
	auditCount    int64
}

var GlobalStore = &Store{}

func (s *Store) GetStats() map[string]int64 {
	return map[string]int64{
		"materials":      atomic.LoadInt64(&s.materialCount),
		"suppliers":      atomic.LoadInt64(&s.supplierCount),
		"inspectors":     atomic.LoadInt64(&s.inspectorCount),
		"lots":           atomic.LoadInt64(&s.lotCount),
		"standards":      atomic.LoadInt64(&s.standardCount),
		"records":        atomic.LoadInt64(&s.recordCount),
		"nonconformances": atomic.LoadInt64(&s.ncCount),
		"return_orders":  atomic.LoadInt64(&s.returnCount),
		"audit_logs":     atomic.LoadInt64(&s.auditCount),
	}
}

func (s *Store) SaveMaterial(m *model.Material) {
	s.Materials.Store(m.ID, m)
	atomic.AddInt64(&s.materialCount, 1)
}

func (s *Store) GetMaterial(id string) (*model.Material, bool) {
	val, ok := s.Materials.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Material), true
}

func (s *Store) DeleteMaterial(id string) {
	if _, ok := s.Materials.LoadAndDelete(id); ok {
		atomic.AddInt64(&s.materialCount, -1)
	}
}

func (s *Store) ListMaterials() []*model.Material {
	result := make([]*model.Material, 0, atomic.LoadInt64(&s.materialCount))
	s.Materials.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Material))
		return true
	})
	return result
}

func (s *Store) SaveSupplier(sup *model.Supplier) {
	s.Suppliers.Store(sup.ID, sup)
	atomic.AddInt64(&s.supplierCount, 1)
}

func (s *Store) GetSupplier(id string) (*model.Supplier, bool) {
	val, ok := s.Suppliers.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Supplier), true
}

func (s *Store) ListSuppliers() []*model.Supplier {
	result := make([]*model.Supplier, 0, atomic.LoadInt64(&s.supplierCount))
	s.Suppliers.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Supplier))
		return true
	})
	return result
}

func (s *Store) SaveInspector(ins *model.Inspector) {
	s.Inspectors.Store(ins.ID, ins)
	atomic.AddInt64(&s.inspectorCount, 1)
}

func (s *Store) GetInspector(id string) (*model.Inspector, bool) {
	val, ok := s.Inspectors.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Inspector), true
}

func (s *Store) ListInspectors() []*model.Inspector {
	result := make([]*model.Inspector, 0, atomic.LoadInt64(&s.inspectorCount))
	s.Inspectors.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Inspector))
		return true
	})
	return result
}

func (s *Store) SaveLot(lot *model.Lot) {
	s.Lots.Store(lot.ID, lot)
	atomic.AddInt64(&s.lotCount, 1)
}

func (s *Store) GetLot(id string) (*model.Lot, bool) {
	val, ok := s.Lots.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Lot), true
}

func (s *Store) GetLotByNo(lotNo string) (*model.Lot, bool) {
	var result *model.Lot
	var found bool
	s.Lots.Range(func(key, value interface{}) bool {
		lot := value.(*model.Lot)
		if lot.LotNo == lotNo {
			result = lot
			found = true
			return false
		}
		return true
	})
	return result, found
}

func (s *Store) ListLots() []*model.Lot {
	result := make([]*model.Lot, 0, atomic.LoadInt64(&s.lotCount))
	s.Lots.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Lot))
		return true
	})
	return result
}

func (s *Store) SaveStandard(std *model.Standard) {
	s.Standards.Store(std.ID, std)
	atomic.AddInt64(&s.standardCount, 1)
}

func (s *Store) GetStandard(id string) (*model.Standard, bool) {
	val, ok := s.Standards.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Standard), true
}

func (s *Store) GetStandardByMaterial(materialID string, typ model.InspectionType) (*model.Standard, bool) {
	var result *model.Standard
	var found bool
	s.Standards.Range(func(key, value interface{}) bool {
		std := value.(*model.Standard)
		if std.MaterialID == materialID && std.Type == typ {
			result = std
			found = true
			return false
		}
		return true
	})
	return result, found
}

func (s *Store) ListStandards() []*model.Standard {
	result := make([]*model.Standard, 0, atomic.LoadInt64(&s.standardCount))
	s.Standards.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Standard))
		return true
	})
	return result
}

func (s *Store) SaveRecord(rec *model.InspectionRecord) {
	s.Records.Store(rec.ID, rec)
	atomic.AddInt64(&s.recordCount, 1)
}

func (s *Store) UpdateRecord(rec *model.InspectionRecord) {
	s.Records.Store(rec.ID, rec)
}

func (s *Store) GetRecord(id string) (*model.InspectionRecord, bool) {
	val, ok := s.Records.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.InspectionRecord), true
}

func (s *Store) ListRecords() []*model.InspectionRecord {
	result := make([]*model.InspectionRecord, 0, atomic.LoadInt64(&s.recordCount))
	s.Records.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.InspectionRecord))
		return true
	})
	return result
}

func (s *Store) GetRecordsByLot(lotID string) []*model.InspectionRecord {
	var result []*model.InspectionRecord
	s.Records.Range(func(key, value interface{}) bool {
		rec := value.(*model.InspectionRecord)
		if rec.LotID == lotID {
			result = append(result, rec)
		}
		return true
	})
	return result
}

func (s *Store) GetRecordsByInspector(inspectorID string) []*model.InspectionRecord {
	var result []*model.InspectionRecord
	s.Records.Range(func(key, value interface{}) bool {
		rec := value.(*model.InspectionRecord)
		if rec.InspectorID == inspectorID && (rec.Status == model.StatusDraft || rec.Status == model.StatusSubmitted) {
			result = append(result, rec)
		}
		return true
	})
	return result
}

func (s *Store) SaveNonconformance(nc *model.Nonconformance) {
	s.Nonconformances.Store(nc.ID, nc)
	atomic.AddInt64(&s.ncCount, 1)
}

func (s *Store) GetNonconformance(id string) (*model.Nonconformance, bool) {
	val, ok := s.Nonconformances.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Nonconformance), true
}

func (s *Store) ListNonconformances() []*model.Nonconformance {
	result := make([]*model.Nonconformance, 0, atomic.LoadInt64(&s.ncCount))
	s.Nonconformances.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.Nonconformance))
		return true
	})
	return result
}

func (s *Store) SaveReturnOrder(ro *model.ReturnOrder) {
	s.ReturnOrders.Store(ro.ID, ro)
	atomic.AddInt64(&s.returnCount, 1)
}

func (s *Store) GetReturnOrder(id string) (*model.ReturnOrder, bool) {
	val, ok := s.ReturnOrders.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.ReturnOrder), true
}

func (s *Store) ListReturnOrders() []*model.ReturnOrder {
	result := make([]*model.ReturnOrder, 0, atomic.LoadInt64(&s.returnCount))
	s.ReturnOrders.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.ReturnOrder))
		return true
	})
	return result
}

func (s *Store) SaveAuditLog(log *model.AuditLog) {
	s.AuditLogs.Store(log.ID, log)
	atomic.AddInt64(&s.auditCount, 1)
}

func (s *Store) GetAuditLogsByRecord(recordID string) []*model.AuditLog {
	var result []*model.AuditLog
	s.AuditLogs.Range(func(key, value interface{}) bool {
		log := value.(*model.AuditLog)
		if log.RecordID == recordID {
			result = append(result, log)
		}
		return true
	})
	return result
}

func (s *Store) SaveApproval(app *model.Approval) {
	s.Approvals.Store(app.ID, app)
}

func (s *Store) GetApproval(id string) (*model.Approval, bool) {
	val, ok := s.Approvals.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Approval), true
}

func (s *Store) UpdateApproval(app *model.Approval) {
	s.Approvals.Store(app.ID, app)
}

func (s *Store) SaveDelegation(d *model.Delegation) {
	s.Delegations.Store(d.ID, d)
}

func (s *Store) GetDelegation(id string) (*model.Delegation, bool) {
	val, ok := s.Delegations.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Delegation), true
}

func (s *Store) GetActiveDelegation(delegatorID string) (*model.Delegation, bool) {
	var result *model.Delegation
	var found bool
	now := timeNow()
	s.Delegations.Range(func(key, value interface{}) bool {
		d := value.(*model.Delegation)
		if d.DelegatorID == delegatorID && d.IsActive && d.StartDate.Before(now) && d.EndDate.After(now) {
			result = d
			found = true
			return false
		}
		return true
	})
	return result, found
}

func (s *Store) SaveSupplierScore(score *model.SupplierScore) {
	s.SupplierScores.Store(score.ID, score)
}

func (s *Store) GetSupplierScores(supplierID string) []*model.SupplierScore {
	var result []*model.SupplierScore
	s.SupplierScores.Range(func(key, value interface{}) bool {
		score := value.(*model.SupplierScore)
		if score.SupplierID == supplierID {
			result = append(result, score)
		}
		return true
	})
	return result
}

func (s *Store) SaveSPCChart(chart *model.SPCControlChart) {
	s.SPCCharts.Store(chart.ID, chart)
}

func (s *Store) GetSPCChart(materialID, itemID string) (*model.SPCControlChart, bool) {
	var result *model.SPCControlChart
	var found bool
	s.SPCCharts.Range(func(key, value interface{}) bool {
		chart := value.(*model.SPCControlChart)
		if chart.MaterialID == materialID && chart.ItemID == itemID {
			result = chart
			found = true
			return false
		}
		return true
	})
	return result, found
}

func (s *Store) ListSPCCharts() []*model.SPCControlChart {
	result := make([]*model.SPCControlChart, 0)
	s.SPCCharts.Range(func(key, value interface{}) bool {
		result = append(result, value.(*model.SPCControlChart))
		return true
	})
	return result
}
