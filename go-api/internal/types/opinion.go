package types

import "time"

type Opinion struct {
	Date         *time.Time   `json:"date" db:"date"`
	Exchange     Exchange     `json:"exchange" db:"exchange" validate:"required,is-exchange"`
	ExpertId     int          `json:"expertId" db:"expert_id" validate:"required"`
	Id           int          `json:"id" db:"id"`
	SecurityId   int          `json:"securityId" db:"security_id"`
	SecurityType SecurityType `json:"securityType" db:"security_type" validate:"required,securityType"`
	SourceLink   *string      `json:"sourceLink" db:"source_link"`
	TargetPrice  *float64     `json:"targetPrice" db:"target_price"`
	Type         OpinionType  `json:"type" db:"type" validate:"required,opinionType"`
	UserId       int          `json:"userId" db:"user_id"`
}

type OpinionType string

const (
	Flat      OpinionType = "Flat"
	General   OpinionType = "General"
	Growth    OpinionType = "Growth"
	Reduction OpinionType = "Reduction"
)

func (e OpinionType) Validate() bool {
	switch e {
	case Flat:
	case General:
	case Growth:
	case Reduction:
	default:
		return false
	}
	return true
}
