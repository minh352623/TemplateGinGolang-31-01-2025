package inittiallize

// func TakeInterest(key string) {
// 	walletRepo := repo.NewWalletRepository()
// 	if walletRepo == nil {
// 		global.Logger.Error("Failed to initialize WalletRepository")
// 		return
// 	}
// 	userRepo := repo.NewUserRepository()
// 	if userRepo == nil {
// 		global.Logger.Error("Failed to initialize UserRepository")
// 		return
// 	}
// 	transactionRepo := repo.NewTransactionRepository()
// 	if transactionRepo == nil {
// 		global.Logger.Error("Failed to initialize TransactionRepository")
// 		return
// 	}
// 	// walletService := service.NewWalletService(walletRepo, userRepo, transactionRepo)
// 	// all users
// 	users, err := userRepo.GetAllUsers()
// 	if err != nil {
// 		global.Logger.Error("Failed to get all users", zap.Error(err))
// 		return
// 	}
// 	for _, user := range users {
// 		//send message queue

// 		data := map[string]interface{}{
// 			"userId":      user.ID,
// 			"providerKey": key,
// 			"type":        consts.TransactionTypeTakeInterest,
// 		}
// 		fmt.Println("data", data)

// 		_, err = messaging.PublishMessage(context.Background(), data, messaging.PublishOptions{
// 			Exchange:     consts.HashedExchangeName,
// 			RoutingKey:   string(messaging.EventTransaction),
// 			HashKey:      user.ID,
// 			WaitResponse: false,
// 		})
// 		if err != nil {
// 			global.Logger.Error("Failed to publish message", zap.Error(err))
// 		}
// 		// walletService.TakeInterest(user.ID)
// 	}
// }

func InitCronJob() {
	// c := cron.New()
	// if global.Config.Cronjob.CronExecuteInterest == "" {
	// 	global.Logger.Error("CronjobSetting is not initialized")
	// 	return
	// }
	// walletIntegrationRepo := repo.NewWalletIntegrationRepository()
	// if walletIntegrationRepo == nil {
	// 	global.Logger.Error("Failed to initialize WalletIntegrationRepository")
	// 	return
	// }
	// walletIntegrations, err := walletIntegrationRepo.GetAllWalletIntegration()
	// if err != nil {
	// 	global.Logger.Error("Failed to get all wallet integrations", zap.Error(err))
	// 	return
	// }
	// for _, walletIntegration := range walletIntegrations {
	// 	if walletIntegration.IsAutoTakeProfit {
	// 		global.Logger.Info("InitCronJob", zap.String("key", walletIntegration.Cronjob))
	// 		_, err = c.AddFunc(walletIntegration.Cronjob, func() {
	// 			TakeInterest(walletIntegration.Key)
	// 		})
	// 		if err != nil {
	// 			global.Logger.Error("InitCronJob", zap.Error(err))
	// 		}
	// 	}
	// }

	// c.Start()
	// fmt.Println("CronJob started")
}
