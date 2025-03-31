package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// MessageHandler is a function that processes a message
type MessageHandler func(ctx context.Context, delivery amqp.Delivery) error

// Consumer is a RabbitMQ message consumer
type Consumer struct {
	rabbitMQ *RabbitMQ
	logger   *zap.Logger
	// Track active workers for graceful shutdown
	activeWorkers sync.WaitGroup
}

// NewConsumer creates a new Consumer
func NewConsumer(rabbitMQ *RabbitMQ, logger *zap.Logger) *Consumer {
	return &Consumer{
		rabbitMQ: rabbitMQ,
		logger:   logger,
	}
}

// Consume starts consuming messages from a queue
func (c *Consumer) Consume(ctx context.Context, queueName, consumerName string, handler MessageHandler) error {
	// Create a context that we can cancel to stop the consumer
	consumerCtx, cancel := context.WithCancel(ctx)

	// Start the consumer
	deliveries, err := c.rabbitMQ.Consume(
		queueName,
		consumerName,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
	)
	if err != nil {
		cancel() // Clean up the context
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Start the main consumer loop
	go func() {
		defer cancel() // Ensure context is cancelled when we exit

		c.logger.Info("Starting consumer",
			zap.String("queue", queueName),
			zap.String("consumer", consumerName))

		for {
			select {
			case <-consumerCtx.Done():
				c.logger.Info("Consumer context cancelled, stopping consumption",
					zap.String("consumer", consumerName))

				// Wait for all workers to finish
				c.activeWorkers.Wait()
				return

			case delivery, ok := <-deliveries:
				if !ok {
					c.logger.Error("Delivery channel closed",
						zap.String("consumer", consumerName))

					// Try to reconnect
					if err := c.rabbitMQ.Reconnect(); err != nil {
						c.logger.Error("Failed to reconnect",
							zap.Error(err),
							zap.String("consumer", consumerName))

						// Wait for all workers to finish
						c.activeWorkers.Wait()
						return
					}

					// Re-register consumer
					newDeliveries, err := c.rabbitMQ.Consume(
						queueName,
						consumerName,
						false, // auto-ack
						false, // exclusive
						false, // no-local
						false, // no-wait
					)
					if err != nil {
						c.logger.Error("Failed to re-register consumer",
							zap.Error(err),
							zap.String("consumer", consumerName))

						// Wait for all workers to finish
						c.activeWorkers.Wait()
						return
					}

					deliveries = newDeliveries
					continue
				}

				// Process the message in a separate goroutine for concurrent processing
				c.activeWorkers.Add(1)
				go func(delivery amqp.Delivery) {
					defer c.activeWorkers.Done()

					// Process the message
					err := handler(consumerCtx, delivery)
					if err != nil {
						c.logger.Error("Error processing message",
							zap.Error(err),
							zap.String("queue", queueName),
							zap.String("consumer", consumerName),
							zap.String("delivery_tag", fmt.Sprintf("%d", delivery.DeliveryTag)),
						)
						// Nack the message to requeue it
						if err := delivery.Nack(false, true); err != nil {
							c.logger.Error("Failed to nack message",
								zap.Error(err),
								zap.String("consumer", consumerName))
						}
					} else {
						// Ack the message
						if err := delivery.Ack(false); err != nil {
							c.logger.Error("Failed to ack message",
								zap.Error(err),
								zap.String("consumer", consumerName))
						}
					}
				}(delivery)
			}
		}
	}()

	return nil
}

// UnmarshalJSON unmarshals a JSON message
func (c *Consumer) UnmarshalJSON(delivery amqp.Delivery, v interface{}) error {
	return json.Unmarshal(delivery.Body, v)
}
