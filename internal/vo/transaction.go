package vo

type GetTransactionByUserIDAndPlatformRequest struct {
	UserID          string `json:"userID"`
	Platform        string `json:"platform"`
	FromDate        string `json:"fromDate"`
	ToDate          string `json:"toDate"`
	TransactionType string `json:"transactionType"`
	Page            int    `json:"page"`
	Limit           int    `json:"limit"`
}
type GetTransactionByCodeRequest struct {
	Code string `json:"code"`
}
