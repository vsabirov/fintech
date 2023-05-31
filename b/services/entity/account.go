package entity

type Account struct {
	ID    string  `db:"id"`
	Total float64 `db:"total"`
}
