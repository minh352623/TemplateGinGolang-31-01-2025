package consts

// define enum
var (
	TransactionTypeDeposit       = "deposit"
	TransactionTypeWithdrawn     = "withdrawn"
	TransactionTypeTakeInterest  = "take-interest"
	TransactionTypeInvestment    = "investment"
	TransactionTypeChargeFee     = "charge-fee"
	TransactionTypeAll           = "all"
	TransactionTypeClaimInterest = "claim-interest"
)

var (
	TransactionIconDeposit       = "icon-deposit"
	TransactionIconWithdraw      = "icon-withdraw"
	TransactionIconTakeInterest  = "icon-take-interest"
	TransactionIconInvestment    = "icon-investment"
	TransactionIconChargeFee     = "icon-charge-fee"
	TransactionIconClaimInterest = "icon-claim-interest"
)

var (
	TransactionStatusPending = "pending"
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)

var (
	HashedExchangeName = "ecom.events.hashed"
)

var (
	WalletIntegrationCurrencyTypeDeposit         = "curency_support_deposit"
	WalletIntegrationCurrencyTypeInputWithdrawn  = "currency_support_input_withdrawn"
	WalletIntegrationCurrencyTypeOutputWithdrawn = "currency_support_output_withdrawn"
)
