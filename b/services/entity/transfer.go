package entity

type Transfer struct {
	ID string `db:"id"`

	Amount float64 `db:"amount"`

	Sender   string `db:"sender"`
	Receiver string `db:"receiver"`
}
