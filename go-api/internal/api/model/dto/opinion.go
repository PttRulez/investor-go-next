package dto

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type Opinion struct {
	Date         *time.Time         `json:"date"`
	Exchange     model.Exchange     `json:"exchange" validate:"required,is-exchange"`
	ExpertId     int                `json:"expertId" validate:"required"`
	Id           int                `json:"id"`
	SecurityId   int                `json:"securityId"`
	SecurityType model.SecurityType `json:"securityType" validate:"required,securityType"`
	SourceLink   *string            `json:"sourceLink"`
	TargetPrice  *float64           `json:"targetPrice"`
	Type         model.OpinionType  `json:"type" validate:"required,opinionType"`
	UserId       int                `json:"userId"`
}
