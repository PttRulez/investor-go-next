package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/types"
	"time"
)

type Opinion struct {
	Date         *time.Time         `json:"date"`
	Exchange     types.Exchange     `json:"exchange" validate:"required,is-exchange"`
	ExpertId     int                `json:"expertId" validate:"required"`
	Id           int                `json:"id"`
	SecurityId   int                `json:"securityId"`
	SecurityType types.SecurityType `json:"securityType" validate:"required,securityType"`
	SourceLink   *string            `json:"sourceLink"`
	TargetPrice  *float64           `json:"targetPrice"`
	Type         entity.OpinionType `json:"type" validate:"required,opinionType"`
	UserId       int                `json:"userId"`
}
