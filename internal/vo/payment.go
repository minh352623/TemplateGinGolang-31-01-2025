package vo

type BalanceChange struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	RateUsd  float64 `json:"rateUsd" binding:"required"`
}
type InvestmentRequest struct {
	UserID          string  `json:"userID" binding:"required"`
	Currency        string  `json:"currency" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	RateUsd         float64 `json:"rateUsd" binding:"required"`
	ProviderKey     string  `json:"providerKey" binding:"required"`
	Platform        string  `json:"platform" binding:"required"`
	WebhookUrl      string  `json:"webhookUrl" binding:"required"`
	TransactionCode string  `json:"transactionCode" binding:"required"`
}
type ChargeFeeRequest struct {
	UserID          string          `json:"userID" binding:"required"`
	ProviderKey     string          `json:"providerKey" binding:"required"`
	Platform        string          `json:"platform" binding:"required"`
	WebhookUrl      string          `json:"webhookUrl" binding:"required"`
	UpdateWallet    []BalanceChange `json:"updateWallet" binding:"required"`
	TransactionCode string          `json:"transactionCode" binding:"required"`
}
