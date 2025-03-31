package vo

import consts "ecom/pkg/const"

type WithdrawRequest struct {
	UserID          string               `json:"userID" binding:"required"`
	Currency        string               `json:"currency" binding:"required"`
	Amount          float64              `json:"amount" binding:"required"`
	RateUsd         float64              `json:"rateUsd" binding:"required"`
	ProviderKey     string               `json:"providerKey" binding:"required"`
	Platform        string               `json:"platform" binding:"required"`
	WebhookUrl      string               `json:"webhookUrl" binding:"required"`
	TransactionCode string               `json:"transactionCode" binding:"required"`
	RateCurrency    consts.CurrencyRates `json:"rateCurrency" binding:"required"`
	ToCurrency      string               `json:"toCurrency" binding:"omitempty"`
}
