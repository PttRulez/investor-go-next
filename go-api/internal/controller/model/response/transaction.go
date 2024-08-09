package response

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"
)

type Transaction struct {
	Amount int                    `json:"amount"`
	Date   time.Time              `json:"date"`
	Id     int                    `json:"id"`
	Type   entity.TransactionType `json:"type"`
}
