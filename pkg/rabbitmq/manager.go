package rabbitmq

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueManager struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	Queues    map[string]amqp.Queue
	Exchanges map[string]string // name -> type
	mu        sync.Mutex
}

type QueueResponse struct {
	CodeResult int
	Data       *[]byte
	Error      string
}

// DeclareExchange safely declares an exchange
func (qm *QueueManager) DeclareExchange(name, kind string) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if err := qm.Channel.ExchangeDeclare(
		name,
		kind,  // direct, fanout, topic, etc.
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return err
	}

	if qm.Exchanges == nil {
		qm.Exchanges = make(map[string]string)
	}
	qm.Exchanges[name] = kind
	return nil
}

// DeclareQueue declares a queue and stores it in the manager
func (qm *QueueManager) DeclareQueue(name string) error {
	// name = test:7 => for and declare queue

	qm.mu.Lock()
	defer qm.mu.Unlock()
	nameNew := strings.Split(name, ":")[0]
	number := strings.Split(name, ":")[1]
	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return err
	}
	for i := 0; i < numberInt; i++ {
		queue, err := qm.Channel.QueueDeclare(
			fmt.Sprintf("%s:%d", nameNew, i),
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,
		)
		if err != nil {
			return err
		}

		if qm.Queues == nil {
			qm.Queues = make(map[string]amqp.Queue)
		}
		qm.Queues[name] = queue
	}
	return nil
}

// BindQueue binds a queue to an exchange with a routing key
func (qm *QueueManager) BindQueue(queueName, exchangeName string) error {
	number := strings.Split(queueName, ":")[1]
	name := strings.Split(queueName, ":")[0]
	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return err
	}
	for i := 0; i < numberInt; i++ {
		queue := fmt.Sprintf("%s:%d", name, i)
		fmt.Println("queue", queue)
		err := qm.Channel.QueueBind(
			queue,
			"1",
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// PublishToExchange sends a message to an exchange with a routing key
func (qm *QueueManager) PublishToExchange(exchange, routingKey, body string) error {
	fmt.Println("PublishToExchange", exchange, routingKey, body)
	return qm.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}

// PublishToExchangeAndWait publishes a message to an exchange and waits for the response
func (qm *QueueManager) PublishToExchangeAndWait(exchange, routingKey, body string, timeout time.Duration) (string, error) {
	// Step 1: Create temporary reply queue
	replyQueue, err := qm.Channel.QueueDeclare(
		"",    // name (let RabbitMQ generate one)
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("declare reply queue: %w", err)
	}

	// Step 2: Start consuming messages from that reply queue
	msgs, err := qm.Channel.Consume(
		replyQueue.Name,
		"",    // consumer tag
		true,  // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("consume reply queue: %w", err)
	}

	// Step 3: Generate unique correlation ID
	corrID := uuid.NewString()

	// Step 4: Publish message to the exchange with replyTo & correlation ID
	err = qm.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
			Body:          []byte(body),
		},
	)
	if err != nil {
		return "", fmt.Errorf("publish failed: %w", err)
	}

	// Step 5: Wait for the matching response
	timeoutChan := time.After(timeout)
	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == corrID {
				return string(msg.Body), nil
			}
		case <-timeoutChan:
			return "", fmt.Errorf("timeout waiting for response")
		}
	}
}

// Consume listens to a queue and processes messages with a handler
func (qm *QueueManager) Consume(queueName string, handler func(msg amqp.Delivery)) error {
	msgs, err := qm.Channel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	return nil
}
