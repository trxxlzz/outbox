package model

import "time"

type Transaction struct {
	ID        string
	UserID    string
	Amount    float64
	Currency  string
	Status    string
	Timestamp time.Time
}
