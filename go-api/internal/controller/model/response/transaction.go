package response

import (
	"time"

	"github.com/pttrulez/investor-go/internal/entity"
)

type Transaction struct {
	Amount int                    `json:"amount"`
	Date   time.Time              `json:"date"`
	ID     int                    `json:"id"`
	Type   entity.TransactionType `json:"type"`
}
