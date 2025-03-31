package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer is a RabbitMQ message producer
type Producer struct {
	rabbitMQ *RabbitMQ
}

// NewProducer creates a new Producer
func NewProducer(rabbitMQ *RabbitMQ) *Producer {
	return &Producer{
		rabbitMQ: rabbitMQ,
	}
}

// PublishRaw publishes a raw message to an exchange
func (p *Producer) PublishRaw(ctx context.Context, exchange, routingKey string, contentType string, data []byte) error {
	msg := amqp.Publishing{
		ContentType:  contentType,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         data,
	}

	return p.rabbitMQ.PublishMessage(ctx, exchange, routingKey, false, false, msg)
}

func (p *Producer) PublishWithPriority(ctx context.Context, exchange, routingKey string, data []byte, priority uint8) error {
	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         data,
		Priority:     priority, // Set message priority
	}

	return p.rabbitMQ.PublishMessage(ctx, exchange, routingKey, false, false, msg)
}

func (p *Producer) PublishHashedMessage(ctx context.Context, exchange, routingKey string, data []byte, hashKey string) error {
	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         data,
		Headers: amqp.Table{
			"hash-on": hashKey, // Messages with same hashKey go to same consumer
		},
	}

	return p.rabbitMQ.PublishMessage(ctx, exchange, routingKey, false, false, msg)
}

// PublishHashedMessageWithResponse publishes a message and waits for response
func (p *Producer) PublishHashedMessageWithResponse(ctx context.Context, exchange, routingKey string, data []byte, hashKey string, timeout ...time.Duration) ([]byte, error) {

	// Define the default timeout as a constant at the package level
	const DefaultResponseTimeout = 5 * time.Minute
	// Use default timeout if not specified
	responseTimeout := DefaultResponseTimeout
	if len(timeout) > 0 && timeout[0] > 0 {
		responseTimeout = timeout[0]
	}

	// Create a unique correlation ID
	correlationID := uuid.New().String()

	// Create a channel for cleanup
	ch, err := p.rabbitMQ.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// Create a temporary response queue
	responseQueue, err := ch.QueueDeclare(
		"",    // let RabbitMQ generate a name
		false, // not durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create response queue: %w", err)
	}

	// Create a channel to receive the response
	responses, err := ch.Consume(
		responseQueue.Name,
		"",    // consumer
		true,  // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create response consumer: %w", err)
	}

	// Publish message with correlation ID and reply-to properties
	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         data,
		Headers: amqp.Table{
			"hash-on": hashKey,
		},
		CorrelationId: correlationID,
		ReplyTo:       responseQueue.Name,
	}

	// Publish the message
	err = p.rabbitMQ.PublishMessage(ctx, exchange, routingKey, false, false, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for response with timeout
	select {
	case delivery := <-responses:
		if delivery.CorrelationId == correlationID {
			return delivery.Body, nil
		}
		return nil, fmt.Errorf("received response with wrong correlation ID")
	case <-time.After(responseTimeout):
		return nil, fmt.Errorf("timeout waiting for response after %v", responseTimeout)
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while waiting for response: %w", ctx.Err())
	}
}
