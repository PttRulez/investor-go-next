package entity

import (
	"github.com/pttrulez/investor-go/internal/types"
	"time"
)

type Opinion struct {
	Date         *time.Time
	Exchange     types.Exchange
	ExpertId     int
	Id           int
	SecurityId   int
	SecurityType types.SecurityType
	SourceLink   *string
	TargetPrice  *float64
	Type         OpinionType
	UserId       int
}

type OpinionType string

const (
	Flat      OpinionType = "FLAT"
	General   OpinionType = "GENERAL"
	Growth    OpinionType = "GROWTH"
	Reduction OpinionType = "REDUCTION"
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
