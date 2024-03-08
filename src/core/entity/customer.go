package entity

type Customer struct {
	ID           string        `json:"id"`
	Limit        int           `json:"limite"`
	Transactions []Transaction `json:"transactions"`
}
