package messaging

import (
	"context"
	"ecom/global"
	consts "ecom/pkg/const"
	"ecom/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// EventHandler is a function that processes an event and returns a response
type EventHandler func(ctx context.Context, event *Event) (interface{}, error)

// MessageCallback is a function called when a message is received or processed
type MessageCallback func()

// patternHandler holds a pattern and its corresponding handler
type patternHandler struct {
	pattern string
	handler EventHandler
}

// ConsumerService manages message consumers
type ConsumerService struct {
	consumer        *rabbitmq.Consumer
	logger          *zap.Logger
	eventHandlers   map[EventType]EventHandler
	patternHandlers []patternHandler
}

// NewConsumerService creates a new consumer service
func NewConsumerService() *ConsumerService {
	consumer := rabbitmq.NewConsumer(global.RabbitMQ, global.Logger.GetZapLogger())
	return &ConsumerService{
		consumer:        consumer,
		logger:          global.Logger.GetZapLogger(),
		eventHandlers:   make(map[EventType]EventHandler),
		patternHandlers: make([]patternHandler, 0),
	}
}

// RegisterEventHandler registers a handler for a specific event type
func (s *ConsumerService) RegisterEventHandler(eventType EventType, handler EventHandler) {
	s.eventHandlers[eventType] = handler
}

// RegisterEventHandlerWithPattern registers a handler for event types matching a pattern
func (s *ConsumerService) RegisterEventHandlerWithPattern(pattern string, handler EventHandler) {
	// Store the pattern and handler in a map
	s.patternHandlers = append(s.patternHandlers, patternHandler{
		pattern: pattern,
		handler: handler,
	})
}

// StartConsumer starts a consumer for the given queue and binding keys
func (s *ConsumerService) StartConsumer(
	ctx context.Context,
	queueName string,
	bindingKeys []string,
	consumerName string,
	onMessageReceived MessageCallback,
	onMessageProcessed MessageCallback,
) error {
	// Declare queue
	queue, err := global.RabbitMQ.DeclareQueue(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to consistent hash exchange with weight
	weight, _ := strconv.Atoi(bindingKeys[0]) // Assuming bindingKeys[0] contains the consumer number
	err = global.RabbitMQ.BindQueueWithWeight(queue.Name, consts.HashedExchangeName, weight)
	if err != nil {
		return fmt.Errorf("failed to bind queue with weight: %w", err)
	}

	// Create a message handler that tracks message processing
	messageHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		// Call the onMessageReceived callback if provided
		if onMessageReceived != nil {
			onMessageReceived()
		}

		// Process the message
		err := s.handleMessage(ctx, delivery)

		// Call the onMessageProcessed callback if provided
		if onMessageProcessed != nil {
			onMessageProcessed()
		}

		return err
	}

	// Start consuming
	err = s.consumer.Consume(ctx, queue.Name, consumerName, messageHandler)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	s.logger.Info("Started consumer",
		zap.String("queue", queueName),
		zap.String("consumer", consumerName),
		zap.Strings("binding_keys", bindingKeys))
	return nil
}

// handleMessage handles incoming messages
func (s *ConsumerService) handleMessage(ctx context.Context, delivery amqp.Delivery) error {
	// Parse the event
	var event Event
	if err := json.Unmarshal(delivery.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Get the event type
	eventType := event.Type

	// Log the received event
	s.logger.Info("Received event",
		zap.String("event_type", string(eventType)),
		zap.Any("payload", event.Payload))

	var response interface{}
	var err error

	// Find and execute the handler
	if handler, exists := s.eventHandlers[eventType]; exists {
		response, err = handler(ctx, &event)
	} else {
		// Check pattern handlers
		eventTypeStr := string(eventType)
		for _, ph := range s.patternHandlers {
			if strings.HasPrefix(eventTypeStr, ph.pattern[:len(ph.pattern)-1]) {
				response, err = ph.handler(ctx, &event)
				break
			}
		}
	}

	// If there's a reply-to queue, send the response
	if delivery.ReplyTo != "" {
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		responseData := map[string]interface{}{
			"success": err == nil,
			"data":    response,
			"error":   errStr,
		}

		responseBytes, err := json.Marshal(responseData)
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		// Publish response using default exchange and reply-to queue as routing key
		err = global.RabbitMQ.PublishMessage(
			ctx,
			"",               // use default exchange
			delivery.ReplyTo, // use reply-to queue as routing key
			false,            // mandatory
			false,            // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				Body:          responseBytes,
				CorrelationId: delivery.CorrelationId, // Important: send back the same correlation ID
			},
		)
		if err != nil {
			return fmt.Errorf("failed to publish response: %w", err)
		}
	}

	return err
}
