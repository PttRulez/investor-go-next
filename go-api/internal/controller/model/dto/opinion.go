package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"
)

type Opinion struct {
	Date         *time.Time          `json:"date"`
	Exchange     entity.Exchange     `json:"exchange" validate:"required,is-exchange"`
	ExpertId     int                 `json:"expertId" validate:"required"`
	Id           int                 `json:"id"`
	SecurityId   int                 `json:"securityId"`
	SecurityType entity.SecurityType `json:"securityType" validate:"required,securityType"`
	SourceLink   *string             `json:"sourceLink"`
	TargetPrice  *float64            `json:"targetPrice"`
	Type         entity.OpinionType  `json:"type" validate:"required,opinionType"`
	UserId       int                 `json:"userId"`
}
