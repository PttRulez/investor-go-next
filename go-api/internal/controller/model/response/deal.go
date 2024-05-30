package response

import (
	"time"
)

type Deal struct {
	Amount int       `json:"amount"`
	Date   time.Time `json:"date"`
	Id     int       `json:"id"`
	Price  float64   `json:"price"`
	Ticker string    `json:"ticker"`
}
