package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"qc-system/config"
	qcron "qc-system/internal/cron"
	"qc-system/internal/handler"
	"qc-system/internal/seed"
	"qc-system/internal/store"
)

func main() {
	cfg := config.Load()

	log.Println("========================================")
	log.Println("QC Quality Management System Starting...")
	log.Println("========================================")

	startTime := time.Now()
	log.Printf("Start time: %s", startTime.Format("2006-01-02 15:04:05"))

	log.Println("Seeding initial data...")
	seed.Seed()
	log.Println("Seed data completed.")

	stats := store.GlobalStore.GetStats()
	log.Println("========================================")
	log.Println("Memory Data Statistics:")
	log.Printf("  Materials:        %d", stats["materials"])
	log.Printf("  Suppliers:        %d", stats["suppliers"])
	log.Printf("  Inspectors:       %d", stats["inspectors"])
	log.Printf("  Lots:             %d", stats["lots"])
	log.Printf("  Standards:        %d", stats["standards"])
	log.Printf("  Records:          %d", stats["records"])
	log.Printf("  Nonconformances:  %d", stats["nonconformances"])
	log.Printf("  Return Orders:    %d", stats["return_orders"])
	log.Printf("  Audit Logs:       %d", stats["audit_logs"])
	log.Println("========================================")
	log.Printf("Sync timestamp: %s", startTime.Format("2006-01-02 15:04:05"))
	log.Println("========================================")

	cronManager := qcron.NewCronManager(cfg)
	cronManager.Start()
	defer cronManager.Stop()

	r := gin.Default()

	r.Use(corsMiddleware())

	setupRoutes(r)

	port := cfg.Port
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	materialHandler := handler.NewMaterialHandler()
	supplierHandler := handler.NewSupplierHandler()
	inspectorHandler := handler.NewInspectorHandler()
	lotHandler := handler.NewLotHandler()
	standardHandler := handler.NewStandardHandler()
	recordHandler := handler.NewRecordHandler()
	ncHandler := handler.NewNCHandler()
	spcHandler := handler.NewSPCHandler()
	statsHandler := handler.NewStatsHandler()
	ocrHandler := handler.NewOCRHandler()

	api := r.Group("/api/v1")
	{
		materials := api.Group("/materials")
		{
			materials.GET("", materialHandler.List)
			materials.POST("", materialHandler.Create)
			materials.GET("/:id", materialHandler.Get)
			materials.PUT("/:id", materialHandler.Update)
			materials.DELETE("/:id", materialHandler.Delete)
		}

		suppliers := api.Group("/suppliers")
		{
			suppliers.GET("", supplierHandler.List)
			suppliers.POST("", supplierHandler.Create)
			suppliers.GET("/focus", supplierHandler.GetFocusSuppliers)
			suppliers.GET("/:id", supplierHandler.Get)
			suppliers.PUT("/:id", supplierHandler.Update)
			suppliers.DELETE("/:id", supplierHandler.Delete)
			suppliers.GET("/:id/scores", supplierHandler.GetScores)
			suppliers.POST("/:id/calculate-score", supplierHandler.CalculateScore)
			suppliers.POST("/calculate-all", supplierHandler.CalculateAll)
		}

		inspectors := api.Group("/inspectors")
		{
			inspectors.GET("", inspectorHandler.List)
			inspectors.POST("", inspectorHandler.Create)
			inspectors.GET("/:id", inspectorHandler.Get)
			inspectors.PUT("/:id", inspectorHandler.Update)
			inspectors.DELETE("/:id", inspectorHandler.Delete)
			inspectors.GET("/:id/records", inspectorHandler.GetRecords)
			inspectors.POST("/:id/delegations", inspectorHandler.CreateDelegation)
			inspectors.GET("/:id/delegations/active", inspectorHandler.GetActiveDelegation)
			inspectors.DELETE("/:id/delegations", inspectorHandler.RevokeDelegation)
			inspectors.POST("/transfer-records", inspectorHandler.TransferRecords)
		}

		lots := api.Group("/lots")
		{
			lots.GET("", lotHandler.List)
			lots.POST("", lotHandler.Create)
			lots.GET("/:id", lotHandler.Get)
			lots.PUT("/:id", lotHandler.Update)
			lots.GET("/:id/records", lotHandler.GetRecords)
			lots.GET("/lot-no/:lot_no", lotHandler.GetByLotNo)
		}

		standards := api.Group("/standards")
		{
			standards.GET("", standardHandler.List)
			standards.POST("", standardHandler.Create)
			standards.GET("/:id", standardHandler.Get)
			standards.PUT("/:id", standardHandler.Update)
			standards.GET("/material/:material_id", standardHandler.GetByMaterial)
			standards.POST("/calculate-sample", standardHandler.CalculateSample)
			standards.GET("/aql/table", standardHandler.GetAQLTable)
		}

		records := api.Group("/records")
		{
			records.GET("", recordHandler.List)
			records.POST("", recordHandler.Create)
			records.GET("/:id", recordHandler.Get)
			records.PUT("/:id", recordHandler.Update)
			records.POST("/:id/submit", recordHandler.Submit)
			records.POST("/:id/approve", recordHandler.Approve)
			records.POST("/:id/reject", recordHandler.Reject)
			records.POST("/:id/close", recordHandler.Close)
			records.POST("/:id/remarks", recordHandler.AddRemark)
			records.POST("/:id/reinspect", recordHandler.Reinspect)
			records.GET("/:id/audit-logs", recordHandler.GetAuditLogs)
			records.GET("/:id/export/pdf", recordHandler.ExportPDF)
			records.GET("/lot/:lot_id", recordHandler.GetByLot)
		}

		nc := api.Group("/nonconformances")
		{
			nc.GET("", ncHandler.List)
			nc.POST("", ncHandler.Create)
			nc.GET("/:id", ncHandler.Get)
			nc.PUT("/:id", ncHandler.Update)
			nc.POST("/:id/close", ncHandler.Close)
			nc.GET("/return-orders", ncHandler.ListReturnOrders)
			nc.POST("/return-orders", ncHandler.CreateReturnOrder)
			nc.GET("/return-orders/:id", ncHandler.GetReturnOrder)
		}

		spc := api.Group("/spc")
		{
			spc.GET("/charts", spcHandler.ListCharts)
			spc.GET("/charts/:material_id", spcHandler.GetChart)
			spc.POST("/charts/:material_id/update", spcHandler.UpdateChart)
			spc.GET("/charts/:material_id/alarms", spcHandler.CheckAlarms)
			spc.POST("/calculate/cl", spcHandler.CalculateControlLimits)
			spc.POST("/calculate/xbar-r", spcHandler.CalculateXBarR)
		}

		stats := api.Group("/stats")
		{
			stats.GET("", statsHandler.GetStats)
			stats.GET("/quality", statsHandler.GetQualityStats)
			stats.GET("/by-material-type", statsHandler.ByMaterialType)
			stats.GET("/by-supplier", statsHandler.BySupplier)
			stats.GET("/by-inspector", statsHandler.ByInspector)
		}

		ocr := api.Group("/ocr")
		{
			ocr.POST("/recognize", ocrHandler.Recognize)
			ocr.POST("/recognize-mark", ocrHandler.RecognizeMarkVerify)
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Operator, X-Role")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
