package services

type TransferRequest struct {
	ID       string  `json:"id" xml:"id"`
	Amount   float64 `json:"amount" xml:"amount"`
	Receiver string  `json:"receiver" xml:"receiver"`
}

func Transfer(request TransferRequest) error {
	return nil
}
