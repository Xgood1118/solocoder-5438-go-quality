package cron

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"qc-system/config"
	"qc-system/internal/model"
	"qc-system/internal/service"
	"qc-system/internal/store"
)

type CronManager struct {
	cfg                *config.Config
	cron               *cron.Cron
	recordService      *service.RecordService
	spcService         *service.SPCService
	supplierScoreService *service.SupplierScoreService
	notificationService *service.NotificationService
}

func NewCronManager(cfg *config.Config) *CronManager {
	return &CronManager{
		cfg:                cfg,
		cron:               cron.New(),
		recordService:      service.NewRecordService(),
		spcService:         service.NewSPCService(),
		supplierScoreService: service.NewSupplierScoreService(),
		notificationService: service.NewNotificationService(cfg),
	}
}

func (m *CronManager) Start() {
	m.cron.AddFunc(m.cfg.InspectionCron, func() { m.PatrolInspectionJob() })
	m.cron.AddFunc(m.cfg.SPCUpdateCron, func() { m.SPCUpdateJob() })
	m.cron.AddFunc(m.cfg.QMSUpdateCron, func() { m.QMSSyncJob() })
	m.cron.AddFunc(m.cfg.SupplierScoreCron, func() { m.SupplierScoreJob() })

	m.cron.Start()
	log.Println("Cron jobs started")
}

func (m *CronManager) Stop() {
	m.cron.Stop()
	log.Println("Cron jobs stopped")
}

func (m *CronManager) PatrolInspectionJob() {
	log.Println("Running patrol inspection job...")

	inspectors := store.GlobalStore.ListInspectors()
	if len(inspectors) == 0 {
		log.Println("No inspectors available for patrol")
		return
	}

	inspectorMap := make(map[string]*model.Inspector)
	for _, ins := range inspectors {
		for _, proc := range ins.Processes {
			if inspectorMap[proc] == nil {
				inspectorMap[proc] = ins
			}
		}
	}

	lots := store.GlobalStore.ListLots()
	allRecords := store.GlobalStore.ListRecords()

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	threeDaysAgo := time.Now().AddDate(0, 0, -3)

	patrolCount := 0

	for _, lot := range lots {
		if lot.ProductionDate.Before(threeDaysAgo) {
			continue
		}

		hasRecentPatrol := false
		for _, rec := range allRecords {
			if rec.LotID == lot.ID && rec.Type == model.TypeIPQC && rec.IPQCType == model.IPQCPatrol && rec.CreatedAt.After(oneHourAgo) {
				hasRecentPatrol = true
				break
			}
		}
		if hasRecentPatrol {
			continue
		}

		material, ok := store.GlobalStore.GetMaterial(lot.MaterialID)
		if !ok {
			continue
		}

		std, ok := store.GlobalStore.GetStandardByMaterial(material.ID, model.TypeIPQC)
		if !ok {
			continue
		}

		items := make([]model.InspectionItemResult, 0, len(std.InspectionLevels))
		for _, item := range std.InspectionLevels {
			items = append(items, model.InspectionItemResult{
				ItemID:    item.ID,
				ItemName:  item.Name,
				Standard:  item.Standard,
				IsNumeric: item.IsNumeric,
				Judgment:  "",
			})
		}

		inspector := inspectorMap[lot.ProductionLine]
		if inspector == nil {
			inspector = inspectors[patrolCount%len(inspectors)]
		}

		record := &model.InspectionRecord{
			ID:              generateID(),
			LotID:           lot.ID,
			LotNo:           lot.LotNo,
			MaterialID:      material.ID,
			SupplierID:      lot.SupplierID,
			Type:            model.TypeIPQC,
			IPQCType:        model.IPQCPatrol,
			InspectorID:     inspector.ID,
			InspectorName:   inspector.Name,
			Status:          model.StatusDraft,
			TotalSampleSize: 5,
			Items:           items,
			Version:         1,
			Remarks:         []model.Remark{},
			CreatedAt:       time.Now(),
		}

		store.GlobalStore.SaveRecord(record)
		patrolCount++
	}

	log.Printf("Patrol inspection job completed, created %d new patrol records", patrolCount)
}

func (m *CronManager) SPCUpdateJob() {
	log.Println("Running SPC update job...")

	charts := store.GlobalStore.ListSPCCharts()
	for _, chart := range charts {
		_, err := m.spcService.UpdateSPCChart(chart.MaterialID, chart.ItemID, chart.ItemName)
		if err != nil {
			log.Printf("Failed to update SPC chart for %s/%s: %v", chart.MaterialID, chart.ItemID, err)
			continue
		}

		updatedChart, ok := m.spcService.GetSPCChart(chart.MaterialID, chart.ItemID)
		if ok {
			alarms := m.spcService.CheckAlarms(updatedChart)
			for _, alarm := range alarms {
				m.notificationService.SendSPCAlert(chart.MaterialID, chart.ItemName, alarm.XBar, chart.UCL, chart.LCL)
			}
		}
	}

	log.Println("SPC update job completed")
}

func (m *CronManager) QMSSyncJob() {
	log.Println("Running QMS sync job...")

	log.Println("QMS sync job completed (simulated)")
}

func (m *CronManager) SupplierScoreJob() {
	log.Println("Running supplier score job...")

	yearMonth := time.Now().AddDate(0, -1, 0).Format("200601")
	scores := m.supplierScoreService.CalculateAllSuppliers(yearMonth)

	log.Printf("Calculated scores for %d suppliers", len(scores))

	for _, score := range scores {
		if score.Level == "red" || score.Level == "yellow" {
			log.Printf("Supplier %s score: %.2f, level: %s", score.SupplierID, score.TotalScore, score.Level)
		}
	}

	log.Println("Supplier score job completed")
}

func generateID() string {
	return "cron-" + time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
