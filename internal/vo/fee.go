package vo

import consts "ecom/pkg/const"

type FeeRequest struct {
	UserID          string               `json:"userID" binding:"required"`
	Currency        string               `json:"currency" binding:"required"`
	Amount          float64              `json:"amount" binding:"required"`
	RateUsd         float64              `json:"rateUsd" binding:"required"`
	ProviderKey     string               `json:"providerKey" binding:"required"`
	Platform        string               `json:"platform" binding:"required"`
	TransactionType string               `json:"transactionType" binding:"required"`
	RateCurrency    consts.CurrencyRates `json:"rateCurrency" binding:"required"`
}
