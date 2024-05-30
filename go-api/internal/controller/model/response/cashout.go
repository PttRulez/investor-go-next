package response

import "time"

type Cashout struct {
	Amount int       `json:"amount"`
	Date   time.Time `json:"date"`
	Id     int       `json:"id"`
	//PortfolioId int       `json:"portfolioId"`
}
