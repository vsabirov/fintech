package services

type Balance struct {
	Total float64 `json:"total" xml:"total"`
}

func GetBalance(accountId string) Balance {
	return Balance{
		Total: 8282.2,
	}
}
