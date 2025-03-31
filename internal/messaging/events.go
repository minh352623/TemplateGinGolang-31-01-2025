package messaging

import (
	"context"
	"ecom/global"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// EventType defines the type of event
type EventType string

// Define event types
const (
	EventUserRegistered EventType = "user.registered"
	EventUserLoggedIn   EventType = "user.logged_in"
	EventWalletCreated  EventType = "wallet.created"
	EventDepositCreated EventType = "deposit.created"
	EventTransaction    EventType = "transaction"
	// Add more event types as needed
)

// Event represents a generic event structure
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// NewEvent creates a new event with the given type and payload
func NewEvent(eventType EventType, payload map[string]interface{}) *Event {
	return &Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// PublishEvent publishes an event to RabbitMQ
func PublishEvent(ctx context.Context, event *Event) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		global.Logger.GetZapLogger().Error("Failed to marshal event",
			zap.Error(err),
			zap.String("event_type", string(event.Type)))
		return err
	}

	// Publish to RabbitMQ
	err = global.RabbitMQProducer.PublishRaw(
		ctx,
		"ecom.events",      // exchange
		string(event.Type), // routing key
		"application/json", // content type
		eventJSON,          // data
	)

	if err != nil {
		global.Logger.GetZapLogger().Error("Failed to publish event",
			zap.Error(err),
			zap.String("event_type", string(event.Type)))
		return err
	}

	global.Logger.GetZapLogger().Info("Published event",
		zap.String("event_type", string(event.Type)),
		zap.Any("payload", event.Payload))
	return nil
}
