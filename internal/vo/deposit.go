package vo

import consts "ecom/pkg/const"

type EncryptedRequest struct {
	Data string `json:"data"`
}
type TestMQRequest struct {
	RoutingKey string `json:"routingKey"`
	MessageID  string `json:"messageID"`
	Email      string `json:"email"`
	UserID     string `json:"userID"`
	HashKey    string `json:"hashKey"`
}
type DepositRequest struct {
	UserID          string               `json:"userID" binding:"required"`
	Currency        string               `json:"currency" binding:"required"`
	RateUsd         float64              `json:"rateUsd" binding:"required"`
	Amount          float64              `json:"amount" binding:"required"`
	ProviderKey     string               `json:"providerKey" binding:"required"`
	TransactionCode string               `json:"transactionCode" binding:"required"`
	Platform        string               `json:"platform" binding:"required"`
	WebhookUrl      string               `json:"webhookUrl" binding:"required"`
	RateCurrency    consts.CurrencyRates `json:"rateCurrency" binding:"required"`
}
