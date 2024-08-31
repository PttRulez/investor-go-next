package converter

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateOpinionRequestToOpinion(ctx context.Context, req api.CreateOpinionRequest) (
	entity.Opinion, error) {
	secType, err := securityType(req.SecurityType)
	if err != nil {
		return entity.Opinion{}, err
	}

	opType, err := opinionType(req.Type)
	if err != nil {
		return entity.Opinion{}, err
	}

	exch, err := exchange(req.Exchange)
	if err != nil {
		return entity.Opinion{}, err
	}

	utils.GetCurrentUserID(ctx)

	return entity.Opinion{
		Date:         entity.Date{Time: req.Date.Time},
		ExpertID:     req.ExpertId,
		Exchange:     exch,
		SecurityID:   req.SecurityId,
		SecurityType: secType,
		SourceLink:   req.SourceLink,
		TargetPrice:  req.TargetPrice,
		Text:         req.Text,
		Ticker:       req.Ticker,
		Type:         opType,
		UserID:       utils.GetCurrentUserID(ctx),
	}, nil
}

func FromOpinionToOpinionResponse(op entity.Opinion) api.OpinionResponse {
	return api.OpinionResponse{
		Date:         openapi_types.Date{Time: op.Date.Time},
		Exchange:     api.Exchange(op.Exchange),
		ExpertId:     op.ExpertID,
		Expert:       FromExpertToExpertResponse(op.Expert),
		Id:           op.ID,
		SecurityId:   op.SecurityID,
		SecurityType: api.SecurityType(op.SecurityType),
		SourceLink:   op.SourceLink,
		TargetPrice:  op.TargetPrice,
		Text:         op.Text,
		Ticker:       op.Ticker,
		Type:         api.OpinionType(op.Type),
	}
}
