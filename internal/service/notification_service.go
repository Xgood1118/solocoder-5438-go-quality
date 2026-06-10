package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"qc-system/config"
)

type NotificationService struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

type AlertMessage struct {
	Type    string            `json:"type"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Level   string            `json:"level"`
	Data    map[string]string `json:"data"`
}

func (s *NotificationService) SendAlert(msg AlertMessage) error {
	var errs []string

	if s.cfg.WebhookURL != "" {
		if err := s.sendWebhook(msg); err != nil {
			errs = append(errs, fmt.Sprintf("webhook: %v", err))
			log.Printf("Webhook send failed: %v", err)
		}
	}

	if s.cfg.SMTPHost != "" {
		if err := s.sendEmail(msg); err != nil {
			errs = append(errs, fmt.Sprintf("email: %v", err))
			log.Printf("Email send failed: %v", err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("some notifications failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (s *NotificationService) sendWebhook(msg AlertMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.cfg.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *NotificationService) sendEmail(msg AlertMessage) error {
	if s.cfg.SMTPHost == "" || s.cfg.SMTPUser == "" {
		return fmt.Errorf("SMTP not configured")
	}

	from := s.cfg.SMTPUser
	to := []string{s.cfg.SMTPUser}
	subject := fmt.Sprintf("[%s] %s", msg.Level, msg.Title)
	body := fmt.Sprintf("Type: %s\r\nLevel: %s\r\n\r\n%s\r\n", msg.Type, msg.Level, msg.Content)

	for k, v := range msg.Data {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, strings.Join(to, ","), subject, body)

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, from, to, []byte(message))
}

func (s *NotificationService) SendSPCAlert(materialID string, itemName string, value float64, ucl, lcl float64) {
	msg := AlertMessage{
		Type:    "spc",
		Title:   "SPC 控制图报警",
		Content: fmt.Sprintf("物料 %s 的 %s 超出控制限，当前值: %.2f, UCL: %.2f, LCL: %.2f", materialID, itemName, value, ucl, lcl),
		Level:   "warning",
		Data: map[string]string{
			"material_id": materialID,
			"item_name":   itemName,
			"value":       fmt.Sprintf("%.2f", value),
			"ucl":         fmt.Sprintf("%.2f", ucl),
			"lcl":         fmt.Sprintf("%.2f", lcl),
		},
	}
	s.SendAlert(msg)
}

func (s *NotificationService) SendQualityAlert(recordID string, lotNo string, judgment string) {
	msg := AlertMessage{
		Type:    "quality",
		Title:   "质量异常报警",
		Content: fmt.Sprintf("批次 %s 检验结果为 %s，记录ID: %s", lotNo, judgment, recordID),
		Level:   "critical",
		Data: map[string]string{
			"record_id": recordID,
			"lot_no":    lotNo,
			"judgment":  judgment,
		},
	}
	s.SendAlert(msg)
}

func (s *NotificationService) SendApprovalNotification(approvalID string, recordID string, currentLevel string) {
	msg := AlertMessage{
		Type:    "approval",
		Title:   "待审批通知",
		Content: fmt.Sprintf("检验记录 %s 待 %s 审批", recordID, currentLevel),
		Level:   "info",
		Data: map[string]string{
			"approval_id": approvalID,
			"record_id":   recordID,
			"level":       currentLevel,
		},
	}
	s.SendAlert(msg)
}
