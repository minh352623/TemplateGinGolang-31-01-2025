package vo

import consts "ecom/pkg/const"

type InterestRequest struct {
	UserID          string               `json:"userID" binding:"required"`
	ProviderKey     string               `json:"providerKey" binding:"required"`
	RateCurrency    consts.CurrencyRates `json:"rateCurrency" binding:"required"`
	WebhookUrl      string               `json:"webhookUrl" binding:"required"`
	Platform        string               `json:"platform" binding:"required"`
	TransactionCode string               `json:"transactionCode" binding:"required"`
}
