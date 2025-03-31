package vo

type GetInfoUserWalletRequest struct {
	UserID      string `json:"userID" binding:"required"`
	ProviderKey string `json:"providerKey" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
}
