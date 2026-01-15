package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"ses-monitoring/internal/config"
	"ses-monitoring/internal/domain/sesevent"
	"ses-monitoring/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SESEvent struct {
	EventType string `json:"eventType"`
	Mail      struct {
		Timestamp     string   `json:"timestamp"`
		MessageID     string   `json:"messageId"`
		Source        string   `json:"source"`
		Destination   []string `json:"destination"`
		CommonHeaders struct {
			Subject string `json:"subject"`
		} `json:"commonHeaders"`
	} `json:"mail"`
	Bounce struct {
		BounceType        string `json:"bounceType"`
		BounceSubType     string `json:"bounceSubType"`
		BouncedRecipients []struct {
			EmailAddress   string `json:"emailAddress"`
			Action         string `json:"action"`
			Status         string `json:"status"`
			DiagnosticCode string `json:"diagnosticCode"`
		} `json:"bouncedRecipients"`
		Timestamp    string `json:"timestamp"`
		ReportingMTA string `json:"reportingMTA"`
	} `json:"bounce"`
	Delivery struct {
		Timestamp            string   `json:"timestamp"`
		ProcessingTimeMillis int      `json:"processingTimeMillis"`
		Recipients           []string `json:"recipients"`
		SmtpResponse         string   `json:"smtpResponse"`
		RemoteMtaIp          string   `json:"remoteMtaIp"`
		ReportingMTA         string   `json:"reportingMTA"`
	} `json:"delivery"`
}

type SNSHandler struct {
	uc              *usecase.SESUsecase
	logBody         bool
	allowedTopicARN string
}

func NewSNSHandler(uc *usecase.SESUsecase, cfg *config.Config) *SNSHandler {
	return &SNSHandler{
		uc:              uc,
		logBody:         cfg.App.LogBody,
		allowedTopicARN: cfg.AWS.SNSTopicARN,
	}
}

func (h *SNSHandler) Handle(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check for SNS Subscription Confirmation
	if typ, ok := payload["Type"].(string); ok {
		if typ == "SubscriptionConfirmation" {
			if subscribeURL, ok := payload["SubscribeURL"].(string); ok {
				log.Printf("SNS Subscription Confirmation received. Auto-confirming...")
				if err := confirmSubscription(subscribeURL); err != nil {
					log.Printf("SNS auto-confirm failed: %v", err)
					c.JSON(http.StatusBadRequest, gin.H{"error": "subscription confirmation failed"})
					return
				}
				log.Printf("SNS subscription confirmed")
				c.JSON(http.StatusOK, gin.H{"status": "subscription confirmed"})
				return
			}
		}
	}

	if h.logBody {
		log.Printf("Received SES payload: %+v", payload)
	}

	if h.allowedTopicARN != "" {
		if topicArn, ok := payload["TopicArn"].(string); !ok || topicArn != h.allowedTopicARN {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid SNS topic"})
			return
		}
	}

	messageStr, ok := payload["Message"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Message field"})
		return
	}

	var sesEvent SESEvent
	if err := json.Unmarshal([]byte(messageStr), &sesEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SES event JSON"})
		return
	}

	// Parse timestamps
	eventTimestamp, err := time.Parse(time.RFC3339, sesEvent.Mail.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail timestamp"})
		return
	}

	// Serialize recipients and tags
	recipientsJSON, _ := json.Marshal(sesEvent.Mail.Destination)
	tagsJSON, _ := json.Marshal(map[string]interface{}{}) // Placeholder, bisa ambil dari payload jika ada

	if len(sesEvent.Mail.Destination) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing destination recipients"})
		return
	}

	// Create event from SES data
	event := &sesevent.Event{
		MessageID:      sesEvent.Mail.MessageID,
		Email:          sesEvent.Mail.Destination[0], // Primary recipient
		Subject:        sesEvent.Mail.CommonHeaders.Subject,
		EventType:      sesEvent.EventType,
		Status:         "SUCCESS",
		Source:         sesEvent.Mail.Source,
		Recipients:     string(recipientsJSON),
		EventTimestamp: eventTimestamp,
		Tags:           string(tagsJSON),
	}

	// Populate based on event type
	switch sesEvent.EventType {
	case "Bounce":
		if len(sesEvent.Bounce.BouncedRecipients) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing bounce recipients"})
			return
		}
		event.Status = "FAILED"
		event.Reason = sesEvent.Bounce.BouncedRecipients[0].DiagnosticCode
		event.BounceType = sesEvent.Bounce.BounceType
		event.BounceSubType = sesEvent.Bounce.BounceSubType
		event.DiagnosticCode = sesEvent.Bounce.BouncedRecipients[0].DiagnosticCode
		event.ReportingMTA = sesEvent.Bounce.ReportingMTA
	case "Delivery":
		event.ProcessingTimeMillis = sesEvent.Delivery.ProcessingTimeMillis
		event.SmtpResponse = sesEvent.Delivery.SmtpResponse
		event.RemoteMtaIp = sesEvent.Delivery.RemoteMtaIp
		event.ReportingMTA = sesEvent.Delivery.ReportingMTA
	}

	err = h.uc.HandleEvent(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func confirmSubscription(subscribeURL string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, subscribeURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
