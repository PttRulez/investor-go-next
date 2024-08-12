package dto

import (
	"time"

	"github.com/pttrulez/investor-go/internal/entity"
)

type Opinion struct {
	Date         *time.Time          `json:"date"`
	Exchange     entity.Exchange     `json:"exchange" validate:"required,is-exchange"`
	ExpertID     int                 `json:"expertId" validate:"required"`
	ID           int                 `json:"id"`
	SecurityID   int                 `json:"securityId"`
	SecurityType entity.SecurityType `json:"securityType" validate:"required,securityType"`
	SourceLink   *string             `json:"sourceLink"`
	TargetPrice  *float64            `json:"targetPrice"`
	Type         entity.OpinionType  `json:"type" validate:"required,opinionType"`
	UserID       int                 `json:"userId"`
}
