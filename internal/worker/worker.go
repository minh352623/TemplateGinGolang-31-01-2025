package worker

// import (
// 	"context"
// 	"ecom/global"
// 	"ecom/internal/messaging"
// 	"ecom/internal/repo"
// 	"ecom/internal/service"
// 	"ecom/internal/utils/interest"
// 	consts "ecom/pkg/const"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"sync"
// 	"sync/atomic"
// 	"syscall"
// 	"time"

// 	"go.uber.org/zap"
// )

// // Worker manages background tasks
// type Worker struct {
// 	consumerService *messaging.ConsumerService
// 	walletService   service.IWalletService
// 	logger          *zap.Logger

// 	// Dynamic scaling properties
// 	activeConsumers     map[string]context.CancelFunc
// 	activeConsumersMu   sync.Mutex
// 	activeMessageCount  int32
// 	minConsumers        int
// 	maxConsumers        int
// 	scalingInterval     time.Duration
// 	messageCountHistory []int32
// 	historySize         int
// }

// // Add at the top with other Worker struct fields
// type TransactionHandler func(ctx context.Context, event *messaging.Event) (interface{}, error)
// type TransactionHandlerMap map[string]TransactionHandler

// // NewWorker creates a new worker
// func NewWorker() *Worker {
// 	// Initialize repositories
// 	walletRepo := repo.NewWalletRepository()
// 	userRepo := repo.NewUserRepository()
// 	transactionRepo := repo.NewTransactionRepository()
// 	walletIntegrationRepo := repo.NewWalletIntegrationRepository()
// 	platformInterestRepo := repo.NewPlatformInterestRepository()
// 	cycleRepo := repo.NewCycleRepository()
// 	transactionTypeRepo := repo.NewTransactionTypeRepository()
// 	walletIntegrationCurrencyRepo := repo.NewWalletIntegrationCurrencyRepository()
// 	projectRepo := repo.NewProjectRepository()
// 	settingService := service.NewSettingService(walletIntegrationRepo, platformInterestRepo, cycleRepo, transactionTypeRepo, walletIntegrationCurrencyRepo, projectRepo)
// 	walletService := service.NewWalletService(walletRepo, userRepo, transactionRepo, settingService)

// 	return &Worker{
// 		consumerService:     messaging.NewConsumerService(),
// 		walletService:       walletService,
// 		logger:              global.Logger.GetZapLogger(),
// 		activeConsumers:     make(map[string]context.CancelFunc),
// 		minConsumers:        2,                // Minimum number of consumers
// 		maxConsumers:        10,               // Maximum number of consumers
// 		scalingInterval:     30 * time.Second, // Check scaling every 30 seconds
// 		messageCountHistory: make([]int32, 0),
// 		historySize:         5, // Keep track of the last 5 measurements
// 	}
// }

// // Start starts the worker
// func (w *Worker) Start() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Handle graceful shutdown
// 	sigCh := make(chan os.Signal, 1)
// 	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

// 	go func() {
// 		sig := <-sigCh
// 		w.logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
// 		cancel()

// 		// Add a timeout for graceful shutdown
// 		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer shutdownCancel()

// 		// Shutdown sequence
// 		w.shutdownAllConsumers()

// 		// Close RabbitMQ connection
// 		if err := global.RabbitMQ.Close(); err != nil {
// 			w.logger.Error("Error closing RabbitMQ connection", zap.Error(err))
// 		}

// 		// Force exit after timeout
// 		select {
// 		case <-shutdownCtx.Done():
// 			w.logger.Warn("Forced shutdown after timeout")
// 			os.Exit(1)
// 		case <-time.After(100 * time.Millisecond):
// 			os.Exit(0)
// 		}
// 	}()

// 	// Register event handlers
// 	w.registerEventHandlers()

// 	// Start initial consumers
// 	w.startInitialConsumers(ctx)

// 	// Start the auto-scaling goroutine
// 	go w.autoScaleConsumers(ctx)

// 	<-ctx.Done()
// 	w.logger.Info("Worker shutting down")
// }

// // startInitialConsumers starts the initial set of consumers
// func (w *Worker) startInitialConsumers(ctx context.Context) {
// 	// Start multiple consumers for consistent hash exchange
// 	for i := 1; i <= w.minConsumers; i++ {
// 		queueName := fmt.Sprintf("sync_queue_%d", i)
// 		w.startConsumer(ctx, queueName, []string{
// 			fmt.Sprintf("%d", i), // Different weight for each queue
// 		}, fmt.Sprintf("sync_consumer_%d", i))
// 	}

// 	// // Other event consumers
// 	// w.startConsumer(ctx, "wallet_events", []string{
// 	// 	string(messaging.EventWalletCreated) + ".*",
// 	// }, "wallet_events_consumer_1")

// 	// w.startConsumer(ctx, "deposit_events", []string{
// 	// 	string(messaging.EventDepositCreated) + ".*",
// 	// }, "deposit_events_consumer_1")
// }

// // startConsumer starts a new consumer with the given parameters
// func (w *Worker) startConsumer(parentCtx context.Context, queueName string, bindingKeys []string, consumerName string) {
// 	// Create a cancellable context for this consumer
// 	ctx, cancel := context.WithCancel(parentCtx)

// 	// Store the cancel function
// 	w.activeConsumersMu.Lock()
// 	w.activeConsumers[consumerName] = cancel
// 	w.activeConsumersMu.Unlock()

// 	// Set QoS (prefetch) for this channel
// 	if err := global.RabbitMQ.Channel().Qos(
// 		1,     // prefetch count - only get 1 message at a time
// 		0,     // prefetch size
// 		false, // global - false means apply to just this channel
// 	); err != nil {
// 		w.logger.Error("Failed to set QoS",
// 			zap.Error(err),
// 			zap.String("consumer", consumerName))
// 		return
// 	}

// 	// Declare queue
// 	queue, err := global.RabbitMQ.DeclareQueue(queueName, true, false, false, false, nil)
// 	if err != nil {
// 		w.logger.Error("Failed to declare queue",
// 			zap.Error(err),
// 			zap.String("queue", queueName))
// 		return
// 	}

// 	// Bind queue to consistent hash exchange with weight
// 	// Use the consumer number as the weight (1, 2, 3, etc.)
// 	weight, _ := strconv.Atoi(bindingKeys[0]) // Assuming bindingKeys[0] contains the consumer number
// 	err = global.RabbitMQ.BindQueueWithWeight(queue.Name, consts.HashedExchangeName, weight)
// 	if err != nil {
// 		w.logger.Error("Failed to bind queue",
// 			zap.Error(err),
// 			zap.String("queue", queueName),
// 			zap.Int("weight", weight))
// 		return
// 	}

// 	// Start the consumer
// 	if err := w.consumerService.StartConsumer(ctx, queueName, bindingKeys, consumerName, w.onMessageReceived, w.onMessageProcessed); err != nil {
// 		w.logger.Error("Failed to start consumer",
// 			zap.Error(err),
// 			zap.String("consumer", consumerName))

// 		w.activeConsumersMu.Lock()
// 		delete(w.activeConsumers, consumerName)
// 		w.activeConsumersMu.Unlock()
// 	} else {
// 		w.logger.Info("Started consumer",
// 			zap.String("consumer", consumerName),
// 			zap.String("queue", queueName))
// 	}
// }

// // stopConsumer stops a specific consumer
// func (w *Worker) stopConsumer(consumerName string) {
// 	w.activeConsumersMu.Lock()
// 	defer w.activeConsumersMu.Unlock()

// 	if cancel, exists := w.activeConsumers[consumerName]; exists {
// 		cancel() // Cancel the context to stop the consumer
// 		delete(w.activeConsumers, consumerName)
// 		w.logger.Info("Stopped consumer", zap.String("consumer", consumerName))
// 	}
// }

// // shutdownAllConsumers stops all active consumers
// func (w *Worker) shutdownAllConsumers() {
// 	w.activeConsumersMu.Lock()
// 	defer w.activeConsumersMu.Unlock()

// 	// Add a WaitGroup to track consumer shutdowns
// 	var wg sync.WaitGroup
// 	wg.Add(len(w.activeConsumers))

// 	for name, cancel := range w.activeConsumers {
// 		go func(name string, cancel context.CancelFunc) {
// 			defer wg.Done()
// 			cancel()
// 			w.logger.Info("Stopped consumer during shutdown", zap.String("consumer", name))
// 		}(name, cancel)
// 	}

// 	// Wait with timeout
// 	done := make(chan struct{})
// 	go func() {
// 		wg.Wait()
// 		close(done)
// 	}()

// 	select {
// 	case <-done:
// 		w.logger.Info("All consumers stopped successfully")
// 	case <-time.After(5 * time.Second):
// 		w.logger.Warn("Timeout waiting for consumers to stop")
// 	}

// 	w.activeConsumers = make(map[string]context.CancelFunc)
// }

// // onMessageReceived is called when a message is received
// func (w *Worker) onMessageReceived() {
// 	atomic.AddInt32(&w.activeMessageCount, 1)
// }

// // onMessageProcessed is called when a message is processed
// func (w *Worker) onMessageProcessed() {
// 	atomic.AddInt32(&w.activeMessageCount, -1)
// }

// // autoScaleConsumers periodically checks if we need to scale up or down
// func (w *Worker) autoScaleConsumers(ctx context.Context) {
// 	ticker := time.NewTicker(w.scalingInterval)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			w.scaleConsumers(ctx)
// 		}
// 	}
// }

// // scaleConsumers decides whether to scale up or down based on current load
// func (w *Worker) scaleConsumers(ctx context.Context) {
// 	// Get current message count
// 	currentCount := atomic.LoadInt32(&w.activeMessageCount)

// 	// Add to history
// 	w.messageCountHistory = append(w.messageCountHistory, currentCount)
// 	if len(w.messageCountHistory) > w.historySize {
// 		w.messageCountHistory = w.messageCountHistory[1:]
// 	}

// 	// Calculate average message count
// 	var sum int32
// 	for _, count := range w.messageCountHistory {
// 		sum += count
// 	}
// 	avgCount := float64(sum) / float64(len(w.messageCountHistory))

// 	// Get current consumer count
// 	w.activeConsumersMu.Lock()
// 	currentConsumerCount := len(w.activeConsumers)
// 	w.activeConsumersMu.Unlock()

// 	// Log current state
// 	w.logger.Info("Auto-scaling check",
// 		zap.Int32("current_message_count", currentCount),
// 		zap.Float64("average_message_count", avgCount),
// 		zap.Int("current_consumer_count", currentConsumerCount))

// 	// Determine if we need to scale up or down
// 	// Scale up if average message count is more than 2 per consumer
// 	if avgCount > float64(currentConsumerCount*2) && currentConsumerCount < w.maxConsumers {
// 		// Scale up by adding more consumers
// 		newConsumersToAdd := min(w.maxConsumers-currentConsumerCount, 2) // Add up to 2 at a time
// 		w.logger.Info("Scaling up consumers", zap.Int("adding_consumers", newConsumersToAdd))

// 		for i := 1; i <= newConsumersToAdd; i++ {
// 			consumerName := fmt.Sprintf("user_events_consumer_%d", currentConsumerCount+i)
// 			w.startConsumer(ctx, "user_events", []string{
// 				string(messaging.EventUserRegistered) + ".*",
// 				string(messaging.EventUserLoggedIn) + ".*",
// 			}, consumerName)
// 		}
// 	}

// 	// Scale down if average message count is less than 1 per consumer and we have more than minimum
// 	if avgCount < float64(currentConsumerCount) && currentConsumerCount > w.minConsumers {
// 		// Scale down by removing consumers
// 		consumersToRemove := min(currentConsumerCount-w.minConsumers, 1) // Remove 1 at a time
// 		w.logger.Info("Scaling down consumers", zap.Int("removing_consumers", consumersToRemove))

// 		// Find the highest numbered consumers to remove
// 		w.activeConsumersMu.Lock()
// 		var consumersToStop []string
// 		for name := range w.activeConsumers {
// 			if len(consumersToStop) < consumersToRemove && name != "wallet_events_consumer_1" && name != "deposit_events_consumer_1" {
// 				// Don't remove the first two user event consumers or the wallet/deposit consumers
// 				if name != "user_events_consumer_1" && name != "user_events_consumer_2" {
// 					consumersToStop = append(consumersToStop, name)
// 				}
// 			}
// 		}
// 		w.activeConsumersMu.Unlock()

// 		// Stop the selected consumers
// 		for _, name := range consumersToStop {
// 			w.stopConsumer(name)
// 		}
// 	}
// }

// // min returns the minimum of two integers
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// // registerEventHandlers registers handlers for different event types
// func (w *Worker) registerEventHandlers() {
// 	// User events
// 	// w.consumerService.RegisterEventHandler(messaging.EventUserRegistered, w.handleUserRegistered)
// 	// Register handler for custom user.registered.* events
// 	w.consumerService.RegisterEventHandlerWithPattern(string(messaging.EventUserRegistered)+".*", w.handleUserRegistered)
// 	w.consumerService.RegisterEventHandlerWithPattern(string(messaging.EventTransaction), w.handleTransaction)
// }

// // Event handlers
// func (w *Worker) handleUserRegistered(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 	userID, _ := event.Payload["user_id"].(string)
// 	email, _ := event.Payload["email"].(string)
// 	messageID, _ := event.Payload["message_id"].(string)

// 	// Extract custom routing key information
// 	eventType := string(event.Type)

// 	w.logger.Info("Starting to process user registered event",
// 		zap.String("user_id", userID),
// 		zap.String("email", email),
// 		zap.String("message_id", messageID),
// 		zap.String("event_type", eventType),
// 		zap.Time("received_at", event.Timestamp))

// 	// Simulate processing time with sleep
// 	// time.Sleep(10 * time.Second)

// 	w.logger.Info("Finished processing user registered event",
// 		zap.String("user_id", userID),
// 		zap.String("email", email),
// 		zap.String("message_id", messageID),
// 		zap.String("event_type", eventType),
// 		zap.Time("processed_at", time.Now()),
// 		zap.Duration("processing_time", time.Since(event.Timestamp)))

// 	// Return a proper response
// 	response := map[string]interface{}{
// 		"status":       "processed",
// 		"user_id":      userID,
// 		"email":        email,
// 		"message_id":   messageID,
// 		"event_type":   eventType,
// 		"processed_at": time.Now(),
// 	}

// 	return response, nil
// }

// func (w *Worker) handleTransaction(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 	userID, _ := event.Payload["userID"].(string)
// 	transactionType, _ := event.Payload["type"].(string)
// 	fmt.Println("userID", userID)
// 	w.logger.Info("Starting to process transaction event",
// 		zap.String("user_id", userID),
// 		zap.String("transaction_type", transactionType))

// 	handlers := w.getTransactionHandlers(userID)

// 	result, err := w.executeTransactionHandler(ctx, event, transactionType, handlers)
// 	if err != nil {
// 		w.logger.Error("Failed to process transaction",
// 			zap.String("type", transactionType),
// 			zap.Error(err))
// 		return nil, err
// 	}
// 	w.logger.Info("Finished processing transaction event",
// 		zap.String("user_id", userID),
// 		zap.String("transaction_type", transactionType),
// 		zap.Any("result", result))

// 	return result, nil
// }

// func (w *Worker) getTransactionHandlers(userID string) TransactionHandlerMap {
// 	return TransactionHandlerMap{
// 		consts.TransactionTypeTakeInterest:  w.handleTakeInterest(userID),
// 		consts.TransactionTypeInvestment:    w.handleInvestment(userID),
// 		consts.TransactionTypeDeposit:       w.handleDeposit(userID),
// 		consts.TransactionTypeChargeFee:     w.handleChargeFee(userID),
// 		consts.TransactionTypeWithdrawn:     w.handleWithdraw(userID),
// 		consts.TransactionTypeClaimInterest: w.handleClaimInterest(userID),

// 		// Add more handlers here as needed
// 	}
// }

// func (w *Worker) handleTakeInterest(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		return w.processWalletUpdateTakeInterest(event, userID)
// 	}
// }

// func (w *Worker) handleInvestment(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		return w.processWalletUpdate(event, userID)
// 	}
// }

// func (w *Worker) handleDeposit(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		w.logger.Info("handleDeposit", zap.String("user_id", userID), zap.Any("event", event))
// 		return w.processWalletUpdate(event, userID)
// 	}
// }
// func (w *Worker) handleClaimInterest(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		return w.processWalletUpdateClaimInterest(event, userID)
// 	}
// }

// func (w *Worker) handleChargeFee(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		return w.processWalletUpdateMultiple(event, userID)
// 	}
// }

// func (w *Worker) handleWithdraw(userID string) TransactionHandler {
// 	return func(ctx context.Context, event *messaging.Event) (interface{}, error) {
// 		return w.processWalletUpdate(event, userID)
// 	}
// }

// func (w *Worker) processWalletUpdateClaimInterest(event *messaging.Event, userID string) (interface{}, error) {
// 	userId := event.Payload["userId"].(string)
// 	providerKey := event.Payload["providerKey"].(string)
// 	platform := event.Payload["platform"].(string)
// 	return w.walletService.ClaimInterest(userId, providerKey, platform)
// }

// func (w *Worker) processWalletUpdateTakeInterest(event *messaging.Event, userID string) (interface{}, error) {
// 	userId := event.Payload["userId"].(string)
// 	providerKey := event.Payload["providerKey"].(string)
// 	return w.walletService.TakeInterest(userId, providerKey)
// }

// func (w *Worker) processWalletUpdate(event *messaging.Event, userID string) (interface{}, error) {
// 	currency := event.Payload["currency"].(string)
// 	amount := event.Payload["amount"].(float64)
// 	providerKey := event.Payload["providerKey"].(string)
// 	interestSettingBytes := event.Payload["interestSetting"].(string)
// 	decodedData, err := base64.StdEncoding.DecodeString(string(interestSettingBytes))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var interestSetting interest.InterestSetting
// 	err = json.Unmarshal(decodedData, &interestSetting)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.walletService.UpdateWallet(userID, currency, amount, providerKey, interestSetting)
// }

// func (w *Worker) processWalletUpdateMultiple(event *messaging.Event, userID string) (interface{}, error) {
// 	updateWalletBytes := event.Payload["updateWallet"].(string)
// 	decodedDataWallet, err := base64.StdEncoding.DecodeString(string(updateWalletBytes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	var updateWallet []service.BalanceChange
// 	err = json.Unmarshal(decodedDataWallet, &updateWallet)
// 	if err != nil {
// 		return nil, err
// 	}
// 	providerKey := event.Payload["providerKey"].(string)
// 	interestSettingBytes := event.Payload["interestSetting"].(string)
// 	decodedData, err := base64.StdEncoding.DecodeString(string(interestSettingBytes))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var interestSetting interest.InterestSetting
// 	err = json.Unmarshal(decodedData, &interestSetting)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.walletService.UpdateWallets(userID, updateWallet, providerKey, interestSetting)
// }

// func (w *Worker) executeTransactionHandler(ctx context.Context, event *messaging.Event, transactionType string, handlers TransactionHandlerMap) (interface{}, error) {
// 	handler, exists := handlers[transactionType]
// 	if !exists {
// 		return nil, fmt.Errorf("unknown transaction type: %s", transactionType)
// 	}

// 	result, err := handler(ctx, event)
// 	if err != nil {
// 		w.logger.Error("Failed to process transaction",
// 			zap.String("type", transactionType),
// 			zap.Error(err))
// 		return nil, err
// 	}

// 	return result, nil
// }

// // Stop gracefully stops the worker
// func (w *Worker) Stop() {
// 	w.shutdownAllConsumers()
// 	w.logger.Info("Worker stopped")
// }
