package converter

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
)

func FromCreateOpinionRequestToOpinion(ctx context.Context, req contracts.CreateOpinionRequest) (
	domain.Opinion, error) {
	secType, err := securityType(req.SecurityType)
	if err != nil {
		return domain.Opinion{}, err
	}

	opType, err := opinionType(req.Type)
	if err != nil {
		return domain.Opinion{}, err
	}

	exch, err := exchange(req.Exchange)
	if err != nil {
		return domain.Opinion{}, err
	}

	utils.GetCurrentUserID(ctx)

	return domain.Opinion{
		Date:         domain.Date{Time: req.Date.Time},
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

func FromOpinionToOpinionResponse(op domain.Opinion) contracts.OpinionResponse {
	return contracts.OpinionResponse{
		Date:         openapi_types.Date{Time: op.Date.Time},
		Exchange:     contracts.Exchange(op.Exchange),
		ExpertId:     op.ExpertID,
		Expert:       FromExpertToExpertResponse(op.Expert),
		Id:           op.ID,
		SecurityId:   op.SecurityID,
		SecurityType: contracts.SecurityType(op.SecurityType),
		SourceLink:   op.SourceLink,
		TargetPrice:  op.TargetPrice,
		Text:         op.Text,
		Ticker:       op.Ticker,
		Type:         contracts.OpinionType(op.Type),
	}
}
