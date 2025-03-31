package tests

import (
	"context"
	"ecom/global"
	"ecom/internal/messaging"
	"ecom/pkg/logger"
	"ecom/pkg/rabbitmq"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// RabbitMQ connection constants
const (
	RabbitMQHost     = "localhost"
	RabbitMQPort     = 5672
	RabbitMQUser     = "guest"
	RabbitMQPassword = "guest"
	RabbitMQVHost    = ""
)

// TestSetup initializes the test environment
func TestSetup(t *testing.T) {
	// Initialize logger if not already initialized
	if global.Logger == nil {
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			t.Fatalf("Failed to initialize logger: %v", err)
		}
		// Wrap the zap logger in a LoggerZap struct
		global.Logger = &logger.LoggerZap{Logger: zapLogger}
	}

	// Initialize RabbitMQ if not already initialized
	if global.RabbitMQ == nil {
		// Create RabbitMQ config using hardcoded values
		config := rabbitmq.Config{
			Host:     RabbitMQHost,
			Port:     RabbitMQPort,
			User:     RabbitMQUser,
			Password: RabbitMQPassword,
			VHost:    RabbitMQVHost,
		}

		// Create RabbitMQ connection
		rabbitMQ, err := rabbitmq.NewRabbitMQ(config, global.Logger.GetZapLogger())
		if err != nil {
			t.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		global.RabbitMQ = rabbitMQ

		// Initialize RabbitMQ producer
		global.RabbitMQProducer = rabbitmq.NewProducer(global.RabbitMQ)
	}

	// Ensure RabbitMQ is connected
	if !global.RabbitMQ.IsConnected() {
		t.Fatal("RabbitMQ is not connected")
	}

	// Declare the exchange for tests
	err := global.RabbitMQ.DeclareExchange("ecom.events", "topic", true, false)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %v", err)
	}
}

// TestParallelMessageProcessing tests sending and processing messages in parallel
func TestParallelMessageProcessing(t *testing.T) {
	TestSetup(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	testCases := []struct {
		name           string
		messageCount   int
		consumerCount  int
		messageDelay   time.Duration
		processingTime time.Duration
	}{
		{"Parallel Processing", 20, 5, 50 * time.Millisecond, 300 * time.Millisecond},
		{"Stress Test", 50, 10, 10 * time.Millisecond, 100 * time.Millisecond},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			queueName := fmt.Sprintf("test_queue_%s", uuid.New().String())

			var processedCount int32
			var processedMessages sync.Map
			var wg sync.WaitGroup
			var processedMutex sync.Mutex
			processedIds := make(map[string]bool)

			consumerService := messaging.NewConsumerService()

			consumerService.RegisterEventHandlerWithPattern("user.registered.*", func(ctx context.Context, event *messaging.Event) (interface{}, error) {
				messageID, _ := event.Payload["message_id"].(string)
				if messageID == "" {
					t.Logf("Warning: Empty message ID received, ignoring...")
					return nil, nil
				}

				// 🚀 Tránh xử lý message trùng
				processedMutex.Lock()
				if processedIds[messageID] {
					processedMutex.Unlock()
					t.Logf("Duplicate message received: %s (ignored)", messageID)
					return nil, nil
				}
				processedIds[messageID] = true
				processedMutex.Unlock()

				// ✅ Gọi wg.Done() chính xác
				defer func() {
					t.Logf("Calling wg.Done() for message: %s", messageID)
					wg.Done()
				}()

				time.Sleep(tc.processingTime)

				processedMessages.Store(messageID, true)
				atomic.AddInt32(&processedCount, 1)
				t.Logf("Processed message: %s, count: %d", messageID, atomic.LoadInt32(&processedCount))

				return nil, nil
			})

			for i := 1; i <= tc.consumerCount; i++ {
				consumerName := fmt.Sprintf("%s_consumer_%d", queueName, i)

				err := consumerService.StartConsumer(
					ctx,
					queueName,
					[]string{"user.registered.*"},
					consumerName,
					func() { t.Logf("Message received by consumer %s", consumerName) },
					func() {},
				)
				require.NoError(t, err)
				t.Logf("Started consumer: %s", consumerName)
			}

			sentMessages := make([]string, 0, tc.messageCount)
			for i := 0; i < tc.messageCount; i++ {
				messageID := uuid.New().String()
				routingKey := fmt.Sprintf("user.registered.message-%d", i%5)
				payload := map[string]interface{}{
					"message_id": messageID,
					"data":       fmt.Sprintf("Test message %d", i),
					"timestamp":  time.Now().UnixNano(),
				}

				event := messaging.NewEvent(messaging.EventType(routingKey), payload)
				err := messaging.PublishEvent(ctx, event)

				if err == nil {
					wg.Add(1) // ✅ Chỉ tăng nếu message gửi thành công
					sentMessages = append(sentMessages, messageID)
					t.Logf("Published message: %s with routing key: %s", messageID, routingKey)
				}

				time.Sleep(tc.messageDelay)
			}

			waitCh := make(chan struct{})
			go func() {
				wg.Wait()
				close(waitCh)
			}()

			select {
			case <-waitCh:
				t.Logf("All messages processed")
			case <-time.After(2 * time.Minute):
				t.Fatalf("Timeout waiting for messages to be processed")
			}

			var totalProcessed int32
			for _, messageID := range sentMessages {
				if _, ok := processedMessages.Load(messageID); ok {
					totalProcessed++
				}
			}

			assert.Equal(t, int32(tc.messageCount), totalProcessed, "Not all messages were processed")
		})
	}
}

// TestDynamicScaling tests the dynamic scaling of consumers
func TestDynamicScaling(t *testing.T) {
	TestSetup(t)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a unique queue name for this test
	queueName := fmt.Sprintf("test_scaling_queue_%s", uuid.New().String())

	// Create a dynamic consumer manager
	type dynamicConsumerManager struct {
		activeConsumers     map[string]context.CancelFunc
		activeConsumersMu   sync.Mutex
		activeMessageCount  int32
		processingCount     int32
		minConsumers        int
		maxConsumers        int
		scalingInterval     time.Duration
		consumerService     *messaging.ConsumerService
		processedMessages   sync.Map
		messageProcessTimes map[string]time.Duration
		messageProcessMu    sync.Mutex
		processedIds        map[string]bool
		processedIdsMu      sync.Mutex
	}

	manager := &dynamicConsumerManager{
		activeConsumers:     make(map[string]context.CancelFunc),
		minConsumers:        2,
		maxConsumers:        10,
		scalingInterval:     5 * time.Second,
		consumerService:     messaging.NewConsumerService(),
		messageProcessTimes: make(map[string]time.Duration),
		processedIds:        make(map[string]bool),
	}

	// Register event handler
	manager.consumerService.RegisterEventHandlerWithPattern("user.logged_in.*", func(ctx context.Context, event *messaging.Event) (interface{}, error) {
		// Track message start time
		messageID, _ := event.Payload["message_id"].(string)

		// Check if we've already processed this message
		manager.processedIdsMu.Lock()
		if _, exists := manager.processedIds[messageID]; exists {
			manager.processedIdsMu.Unlock()
			t.Logf("Duplicate processing of message: %s (ignored)", messageID)
			return nil, nil
		}
		manager.processedIds[messageID] = true
		manager.processedIdsMu.Unlock()

		startTime := time.Now()

		// Increment processing count
		atomic.AddInt32(&manager.processingCount, 1)

		// Simulate variable processing time based on message content
		processingTime, _ := event.Payload["processing_time"].(float64)
		if processingTime > 0 {
			time.Sleep(time.Duration(processingTime) * time.Millisecond)
		} else {
			// Default processing time
			time.Sleep(200 * time.Millisecond)
		}

		// Mark message as processed
		manager.processedMessages.Store(messageID, true)

		// Store processing time
		manager.messageProcessMu.Lock()
		manager.messageProcessTimes[messageID] = time.Since(startTime)
		manager.messageProcessMu.Unlock()

		// Decrement processing count
		atomic.AddInt32(&manager.processingCount, -1)

		return nil, nil
	})

	// Function to start a consumer
	startConsumer := func(parentCtx context.Context, queueName string, bindingKeys []string, consumerName string) {
		// Create a cancellable context for this consumer
		consumerCtx, cancel := context.WithCancel(parentCtx)

		// Store the cancel function
		manager.activeConsumersMu.Lock()
		manager.activeConsumers[consumerName] = cancel
		manager.activeConsumersMu.Unlock()

		// Start the consumer
		err := manager.consumerService.StartConsumer(
			consumerCtx,
			queueName,
			bindingKeys,
			consumerName,
			func() {
				atomic.AddInt32(&manager.activeMessageCount, 1)
			},
			func() {
				atomic.AddInt32(&manager.activeMessageCount, -1)
			},
		)

		if err != nil {
			t.Logf("Failed to start consumer %s: %v", consumerName, err)

			// Remove from active consumers if failed to start
			manager.activeConsumersMu.Lock()
			delete(manager.activeConsumers, consumerName)
			manager.activeConsumersMu.Unlock()
		} else {
			t.Logf("Started consumer: %s", consumerName)
		}
	}

	// Function to auto-scale consumers
	autoScaleConsumers := func(ctx context.Context) {
		ticker := time.NewTicker(manager.scalingInterval)
		defer ticker.Stop()

		// Helper function for min
		min := func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get current message and processing counts
				currentMessageCount := atomic.LoadInt32(&manager.activeMessageCount)
				currentProcessingCount := atomic.LoadInt32(&manager.processingCount)

				// Get current consumer count
				manager.activeConsumersMu.Lock()
				currentConsumerCount := len(manager.activeConsumers)
				manager.activeConsumersMu.Unlock()

				t.Logf("Auto-scaling check - Messages: %d, Processing: %d, Consumers: %d",
					currentMessageCount, currentProcessingCount, currentConsumerCount)

				// Scale up if we have more messages than consumers can handle
				if currentMessageCount > int32(currentConsumerCount*2) && currentConsumerCount < manager.maxConsumers {
					// Scale up by adding more consumers
					newConsumersToAdd := min(manager.maxConsumers-currentConsumerCount, 2)
					t.Logf("Scaling up by adding %d consumers", newConsumersToAdd)

					for i := 1; i <= newConsumersToAdd; i++ {
						consumerName := fmt.Sprintf("%s_consumer_%d", queueName, currentConsumerCount+i)
						startConsumer(ctx, queueName, []string{"user.logged_in.*"}, consumerName)
					}
				}

				// Scale down if we have fewer messages than consumers
				if currentMessageCount < int32(currentConsumerCount) && currentConsumerCount > manager.minConsumers {
					// Scale down by removing one consumer
					t.Logf("Scaling down by removing 1 consumer")

					// Find a consumer to remove
					manager.activeConsumersMu.Lock()
					var consumerToStop string
					for name := range manager.activeConsumers {
						if name != fmt.Sprintf("%s_consumer_1", queueName) &&
							name != fmt.Sprintf("%s_consumer_2", queueName) {
							consumerToStop = name
							break
						}
					}

					// Stop the consumer
					if consumerToStop != "" {
						cancel := manager.activeConsumers[consumerToStop]
						delete(manager.activeConsumers, consumerToStop)
						manager.activeConsumersMu.Unlock()

						cancel()
						t.Logf("Stopped consumer: %s", consumerToStop)
					} else {
						manager.activeConsumersMu.Unlock()
					}
				}
			}
		}
	}

	// Start the initial consumers
	for i := 1; i <= manager.minConsumers; i++ {
		startConsumer(ctx, queueName, []string{"user.logged_in.*"}, fmt.Sprintf("%s_consumer_%d", queueName, i))
	}

	// Start the auto-scaling goroutine
	go autoScaleConsumers(ctx)

	// Test phases
	phases := []struct {
		name           string
		messageCount   int
		messageDelay   time.Duration
		processingTime time.Duration
		burstFactor    int // 1 = steady, >1 = burst
	}{
		{
			name:           "Steady Low Load",
			messageCount:   20,
			messageDelay:   500 * time.Millisecond,
			processingTime: 200 * time.Millisecond,
			burstFactor:    1,
		},
		{
			name:           "Medium Load",
			messageCount:   50,
			messageDelay:   200 * time.Millisecond,
			processingTime: 300 * time.Millisecond,
			burstFactor:    1,
		},
		{
			name:           "Burst Load",
			messageCount:   30,
			messageDelay:   50 * time.Millisecond,
			processingTime: 500 * time.Millisecond,
			burstFactor:    3,
		},
		{
			name:           "High Sustained Load",
			messageCount:   100,
			messageDelay:   100 * time.Millisecond,
			processingTime: 400 * time.Millisecond,
			burstFactor:    1,
		},
		{
			name:           "Cool Down Period",
			messageCount:   10,
			messageDelay:   1 * time.Second,
			processingTime: 100 * time.Millisecond,
			burstFactor:    1,
		},
	}

	// Run test phases
	sentMessages := make([]string, 0)

	for _, phase := range phases {
		t.Logf("Starting phase: %s", phase.name)

		// Send messages for this phase
		for i := 0; i < phase.messageCount; i++ {
			// Create a unique message ID
			messageID := uuid.New().String()
			sentMessages = append(sentMessages, messageID)

			// Create custom routing key with suffix
			routingKey := fmt.Sprintf("user.logged_in.message-%d", i%5)

			// Determine processing time - add some randomness
			processingTime := float64(phase.processingTime)

			// Create event payload
			payload := map[string]interface{}{
				"message_id":      messageID,
				"data":            fmt.Sprintf("Test message %d - %s", i, phase.name),
				"timestamp":       time.Now().UnixNano(),
				"processing_time": processingTime,
				"phase":           phase.name,
			}

			// Create and publish event
			event := messaging.NewEvent(messaging.EventType(routingKey), payload)
			err := messaging.PublishEvent(ctx, event)
			require.NoError(t, err)

			// For burst factor, send multiple messages at once
			if phase.burstFactor > 1 && i%(phase.burstFactor*2) == 0 {
				// Send a burst of messages
				for j := 0; j < phase.burstFactor-1; j++ {
					burstMessageID := uuid.New().String()
					sentMessages = append(sentMessages, burstMessageID)

					burstPayload := map[string]interface{}{
						"message_id":      burstMessageID,
						"data":            fmt.Sprintf("Burst message %d.%d - %s", i, j, phase.name),
						"timestamp":       time.Now().UnixNano(),
						"processing_time": processingTime * 0.8, // Slightly faster processing for burst messages
						"phase":           phase.name + " (burst)",
					}

					burstEvent := messaging.NewEvent(messaging.EventType(routingKey), burstPayload)
					err = messaging.PublishEvent(ctx, burstEvent)
					require.NoError(t, err)
				}
			}

			// Delay between messages
			time.Sleep(phase.messageDelay)
		}

		// Wait a bit between phases
		time.Sleep(5 * time.Second)
	}

	// Wait for all messages to be processed
	waitForProcessing := func() bool {
		deadline := time.Now().Add(1 * time.Minute)
		for time.Now().Before(deadline) {
			processedCount := 0
			for _, messageID := range sentMessages {
				if _, ok := manager.processedMessages.Load(messageID); ok {
					processedCount++
				}
			}

			if processedCount == len(sentMessages) {
				return true
			}

			t.Logf("Waiting for messages to be processed: %d/%d", processedCount, len(sentMessages))
			time.Sleep(2 * time.Second)
		}
		return false
	}

	// Wait for processing to complete
	allProcessed := waitForProcessing()
	assert.True(t, allProcessed, "Not all messages were processed within the timeout")

	// Calculate statistics
	var totalProcessingTime time.Duration
	var maxProcessingTime time.Duration
	var minProcessingTime = time.Hour

	manager.messageProcessMu.Lock()
	for _, duration := range manager.messageProcessTimes {
		totalProcessingTime += duration
		if duration > maxProcessingTime {
			maxProcessingTime = duration
		}
		if duration < minProcessingTime {
			minProcessingTime = duration
		}
	}
	manager.messageProcessMu.Unlock()

	avgProcessingTime := totalProcessingTime / time.Duration(len(sentMessages))

	t.Logf("Processing statistics:")
	t.Logf("  Total messages: %d", len(sentMessages))
	t.Logf("  Average processing time: %v", avgProcessingTime)
	t.Logf("  Min processing time: %v", minProcessingTime)
	t.Logf("  Max processing time: %v", maxProcessingTime)
}

// TestMessageBatching tests sending and processing messages in batches
func TestMessageBatching(t *testing.T) {
	TestSetup(t)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Create a unique queue name for this test
	queueName := fmt.Sprintf("test_batch_queue_%s", uuid.New().String())

	// Track processed messages
	var processedCount int32
	var processedMessages sync.Map
	var wg sync.WaitGroup
	var processedMutex sync.Mutex
	processedIds := make(map[string]bool)

	// Create consumer service
	consumerService := messaging.NewConsumerService()

	// Register event handler
	consumerService.RegisterEventHandlerWithPattern("wallet.created.*", func(ctx context.Context, event *messaging.Event) (interface{}, error) {
		// Simulate processing time
		time.Sleep(200 * time.Millisecond)

		// Mark message as processed
		messageID, _ := event.Payload["message_id"].(string)

		// Use mutex to safely check and update processed IDs
		processedMutex.Lock()
		if _, exists := processedIds[messageID]; !exists {
			// First time seeing this message
			processedIds[messageID] = true
			processedMutex.Unlock()

			processedMessages.Store(messageID, true)
			atomic.AddInt32(&processedCount, 1)

			t.Logf("Processed message: %s, count: %d", messageID, atomic.LoadInt32(&processedCount))
			wg.Done() // Only call Done() once per message
		} else {
			processedMutex.Unlock()
			t.Logf("Duplicate processing of message: %s (ignored for WaitGroup)", messageID)
		}

		return nil, nil
	})

	// Start multiple consumers
	consumerCount := 5
	for i := 1; i <= consumerCount; i++ {
		consumerName := fmt.Sprintf("%s_consumer_%d", queueName, i)

		// Start consumer with callbacks
		err := consumerService.StartConsumer(
			ctx,
			queueName,
			[]string{"wallet.created.*"},
			consumerName,
			func() {},
			func() {},
		)
		require.NoError(t, err)
		t.Logf("Started consumer: %s", consumerName)
	}

	// Test sending messages in batches
	batchSizes := []int{10, 20, 50, 100}
	totalMessages := 0
	sentMessages := make([]string, 0)

	for _, batchSize := range batchSizes {
		t.Logf("Sending batch of %d messages", batchSize)

		// Send batch of messages
		for i := 0; i < batchSize; i++ {
			// Create a unique message ID
			messageID := uuid.New().String()
			sentMessages = append(sentMessages, messageID)

			// Create custom routing key with suffix
			routingKey := fmt.Sprintf("wallet.created.batch-%d", len(batchSizes))

			// Create event payload
			payload := map[string]interface{}{
				"message_id": messageID,
				"data":       fmt.Sprintf("Batch message %d of size %d", i, batchSize),
				"batch_size": batchSize,
				"timestamp":  time.Now().UnixNano(),
			}

			// Add to wait group to track this message
			wg.Add(1)

			// Create and publish event
			event := messaging.NewEvent(messaging.EventType(routingKey), payload)
			err := messaging.PublishEvent(ctx, event)
			require.NoError(t, err)
		}

		totalMessages += batchSize

		// Wait for batch to be processed with timeout
		batchWaitCh := make(chan struct{})
		go func() {
			wg.Wait()
			close(batchWaitCh)
		}()

		select {
		case <-batchWaitCh:
			t.Logf("Batch of %d messages processed", batchSize)
		case <-time.After(1 * time.Minute):
			t.Fatalf("Timeout waiting for batch of %d messages to be processed", batchSize)
		}

		// Verify all messages in this batch were processed
		processedBatchCount := 0
		for _, messageID := range sentMessages[totalMessages-batchSize:] {
			if _, ok := processedMessages.Load(messageID); ok {
				processedBatchCount++
			}
		}

		assert.Equal(t, batchSize, processedBatchCount, "Not all messages in batch were processed")

		// Wait between batches
		time.Sleep(2 * time.Second)
	}

	// Verify all messages were processed
	assert.Equal(t, int32(totalMessages), atomic.LoadInt32(&processedCount), "Not all messages were processed")
}
