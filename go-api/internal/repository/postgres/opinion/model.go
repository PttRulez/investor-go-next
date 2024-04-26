package opinion

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type Opinion struct {
	Date         *time.Time         `db:"date"`
	Exchange     model.Exchange     `db:"exchange"`
	ExpertId     int                `db:"expert_id" `
	Id           int                `db:"id"`
	SecurityId   int                `db:"security_id"`
	SecurityType model.SecurityType `db:"security_type"`
	SourceLink   *string            `db:"source_link"`
	TargetPrice  *float64           `db:"target_price"`
	Type         model.OpinionType  `db:"type"`
	UserId       int                `db:"user_id"`
}
