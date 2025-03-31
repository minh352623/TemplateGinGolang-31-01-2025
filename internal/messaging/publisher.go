package messaging

import (
	"context"
	"ecom/global"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type PublishOptions struct {
	Exchange     string
	RoutingKey   string
	HashKey      string
	WaitResponse bool
	Timeout      time.Duration
}

// Add a generic response struct
type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Error   string `json:"error"`
}

// PublishMessage publishes a message with optional response waiting
func PublishMessage(ctx context.Context, data interface{}, opts PublishOptions) ([]byte, error) {
	// Create event
	event := NewEvent(EventType(opts.RoutingKey), data.(map[string]interface{}))

	// Marshal event
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// If wait for response
	if opts.WaitResponse {
		return global.RabbitMQProducer.PublishHashedMessageWithResponse(
			ctx,
			opts.Exchange,
			opts.RoutingKey,
			eventJSON,
			opts.HashKey,
			opts.Timeout,
		)
	}

	// Fire and forget
	err = global.RabbitMQProducer.PublishHashedMessage(
		ctx,
		opts.Exchange,
		opts.RoutingKey,
		eventJSON,
		opts.HashKey,
	)
	return nil, err
}

// Modify PublishMessage to use generics
func PublishMessageWithResponse[T any](ctx context.Context, data interface{}, opts PublishOptions) (*T, error) {
	responseBytes, err := PublishMessage(ctx, data, opts)
	if err != nil {
		return nil, err
	}

	var result Response[T]
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, errors.New(result.Error)
	}

	return &result.Data, nil
}
