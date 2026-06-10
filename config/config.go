package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	WebhookURL      string
	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	InspectionCron  string
	SPCUpdateCron   string
	QMSUpdateCron   string
	QMSBaseURL      string
	SupplierScoreCron string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         getEnvInt("SMTP_PORT", 587),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPassword:     getEnv("SMTP_PASSWORD", ""),
		InspectionCron:   getEnv("INSPECTION_CRON", "0 * * * *"),
		SPCUpdateCron:    getEnv("SPC_UPDATE_CRON", "0 */4 * * *"),
		QMSUpdateCron:    getEnv("QMS_UPDATE_CRON", "0 2 * * *"),
		QMSBaseURL:       getEnv("QMS_BASE_URL", "http://qms.internal/api"),
		SupplierScoreCron: getEnv("SUPPLIER_SCORE_CRON", "0 3 1 * *"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
