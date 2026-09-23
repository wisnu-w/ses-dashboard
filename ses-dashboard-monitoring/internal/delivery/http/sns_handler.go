package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"ses-monitoring/internal/config"
	"ses-monitoring/internal/domain/sesevent"
	"ses-monitoring/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SESEvent struct {
	EventType string `json:"eventType"`
	Mail      struct {
		Timestamp     string                 `json:"timestamp"`
		MessageID     string                 `json:"messageId"`
		Source        string                 `json:"source"`
		Destination   []string               `json:"destination"`
		Tags          map[string]interface{} `json:"tags"`
		CommonHeaders struct {
			Subject string   `json:"subject"`
			To      []string `json:"to"`
			Cc      []string `json:"cc"`
			Bcc     []string `json:"bcc"`
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
	Complaint struct {
		ComplainedRecipients []struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"complainedRecipients"`
	} `json:"complaint"`
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

func getDestinationType(email string, headers struct {
	Subject string   `json:"subject"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
}) string {
	for _, to := range headers.To {
		if strings.Contains(to, email) {
			return "To"
		}
	}
	for _, cc := range headers.Cc {
		if strings.Contains(cc, email) {
			return "Cc"
		}
	}
	for _, bcc := range headers.Bcc {
		if strings.Contains(bcc, email) {
			return "Bcc"
		}
	}
	return "To"
}

func (h *SNSHandler) Handle(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if typ, ok := payload["Type"].(string); ok {
		if typ == "SubscriptionConfirmation" {
			if subscribeURL, ok := payload["SubscribeURL"].(string); ok {
				log.Printf("SNS Subscription Confirmation received. SubscribeURL: %s", subscribeURL)
				c.JSON(http.StatusOK, gin.H{"status": "subscription confirmation logged"})
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

	eventTimestamp, err := time.Parse(time.RFC3339, sesEvent.Mail.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail timestamp"})
		return
	}

	if sesEvent.Mail.MessageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing mail message ID"})
		return
	}

	if len(sesEvent.Mail.Destination) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing destination recipients"})
		return
	}

	recipientsJSON, _ := json.Marshal(sesEvent.Mail.Destination)
	tagsJSON, _ := json.Marshal(sesEvent.Mail.Tags)

	baseEvent := sesevent.Event{
		MessageID:      sesEvent.Mail.MessageID,
		Subject:        sesEvent.Mail.CommonHeaders.Subject,
		EventType:      sesEvent.EventType,
		Status:         "SUCCESS",
		Source:         sesEvent.Mail.Source,
		Recipients:     string(recipientsJSON),
		EventTimestamp: eventTimestamp,
		Tags:           string(tagsJSON),
	}

	var eventsToSave []*sesevent.Event

	switch sesEvent.EventType {
	case "Bounce":
		for _, rec := range sesEvent.Bounce.BouncedRecipients {
			e := baseEvent
			e.Email = rec.EmailAddress
			e.Status = "FAILED"
			e.Reason = rec.DiagnosticCode
			e.BounceType = sesEvent.Bounce.BounceType
			e.BounceSubType = sesEvent.Bounce.BounceSubType
			e.DiagnosticCode = rec.DiagnosticCode
			e.ReportingMTA = sesEvent.Bounce.ReportingMTA
			e.DestinationType = getDestinationType(rec.EmailAddress, sesEvent.Mail.CommonHeaders)
			eventsToSave = append(eventsToSave, &e)
		}
	case "Complaint":
		for _, rec := range sesEvent.Complaint.ComplainedRecipients {
			e := baseEvent
			e.Email = rec.EmailAddress
			e.Status = "COMPLAINT"
			e.DestinationType = getDestinationType(rec.EmailAddress, sesEvent.Mail.CommonHeaders)
			eventsToSave = append(eventsToSave, &e)
		}
	case "Delivery":
		for _, rec := range sesEvent.Delivery.Recipients {
			e := baseEvent
			e.Email = rec
			e.ProcessingTimeMillis = sesEvent.Delivery.ProcessingTimeMillis
			e.SmtpResponse = sesEvent.Delivery.SmtpResponse
			e.RemoteMtaIp = sesEvent.Delivery.RemoteMtaIp
			e.ReportingMTA = sesEvent.Delivery.ReportingMTA
			e.DestinationType = getDestinationType(rec, sesEvent.Mail.CommonHeaders)
			eventsToSave = append(eventsToSave, &e)
		}
	case "Send":
		baseEvent.Status = "PENDING"
	case "Open", "Click":
		baseEvent.Status = "SUCCESS"
	default:
		baseEvent.Status = "UNKNOWN"
	}

	if len(eventsToSave) == 0 {
		for _, rec := range sesEvent.Mail.Destination {
			e := baseEvent
			e.Email = rec
			e.DestinationType = getDestinationType(rec, sesEvent.Mail.CommonHeaders)
			eventsToSave = append(eventsToSave, &e)
		}
	}

	for _, e := range eventsToSave {
		if err := h.uc.HandleEvent(c.Request.Context(), e); err != nil {
			log.Printf("failed to persist SES event type=%s message_id=%s: %v", e.EventType, e.MessageID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
