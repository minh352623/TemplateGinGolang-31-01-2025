package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMQ connection manager
type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	uri          string
	exchangeName string
	logger       *zap.Logger
}

// Config holds RabbitMQ connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(config Config, logger *zap.Logger) (*RabbitMQ, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		uri:     uri,
		logger:  logger,
	}, nil
}

// Close closes the connection and channel
func (r *RabbitMQ) Close() error {
	var errs []error

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing channel: %w", err))
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during RabbitMQ shutdown: %v", errs)
	}
	return nil
}

// DeclareExchange declares an exchange
func (r *RabbitMQ) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	r.exchangeName = name
	return r.channel.ExchangeDeclare(
		name,       // name
		kind,       // type
		durable,    // durable
		autoDelete, // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
}

// DeclareQueue declares a queue
func (r *RabbitMQ) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return r.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

// BindQueue binds a queue to an exchange
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// PublishMessage publishes a message to an exchange
func (r *RabbitMQ) PublishMessage(ctx context.Context, exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	return r.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		mandatory,  // mandatory
		immediate,  // immediate
		msg,        // message
	)
}

// Consume consumes messages from a queue
func (r *RabbitMQ) Consume(queueName, consumerName string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,    // queue
		consumerName, // consumer
		autoAck,      // auto-ack
		exclusive,    // exclusive
		noLocal,      // no-local
		noWait,       // no-wait
		nil,          // args
	)
}

// IsConnected checks if connection is established
func (r *RabbitMQ) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}

// Reconnect attempts to reconnect to RabbitMQ
func (r *RabbitMQ) Reconnect() error {
	if r.IsConnected() {
		return nil
	}

	conn, err := amqp.Dial(r.uri)
	if err != nil {
		return fmt.Errorf("failed to reconnect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	// Close old connection and channel if they exist
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

	r.conn = conn
	r.channel = ch
	return nil
}

// DeclareConsistentHashExchange declares a consistent hash exchange
func (r *RabbitMQ) DeclareConsistentHashExchange(name string) error {
	args := amqp.Table{
		"hash-header": "hash-on",
		// "hash-property": "message_id", // Fallback to message_id if header not present
	}

	return r.channel.ExchangeDeclare(
		name,
		"x-consistent-hash",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		args,
	)
}

// BindQueueWithWeight binds a queue to a consistent hash exchange with a weight
func (r *RabbitMQ) BindQueueWithWeight(queueName, exchangeName string, weight int) error {
	args := amqp.Table{
		"hash-header": "hash-on",
	}

	// The routing key for consistent-hash exchanges is the weight as a string
	return r.channel.QueueBind(
		queueName,                 // queue name
		fmt.Sprintf("%d", weight), // routing key (weight)
		exchangeName,              // exchange
		false,                     // no-wait
		args,                      // arguments
	)
}

// Channel returns the AMQP channel
func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.channel
}
