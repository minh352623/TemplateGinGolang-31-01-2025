package consts

type CurrencyRate struct {
	USD float64 `json:"USD"`
}

type CurrencyRates map[string]CurrencyRate
type Result struct {
	Currency       string  `json:"currency"`
	AmountInterest float64 `json:"amountInterest"`
}
