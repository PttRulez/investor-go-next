package converter

import (
	"context"

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
		Type:         opType,
		UserID:       utils.GetCurrentUserID(ctx),
	}, nil
}

// func OpinionFiltersReqToOpinionFilters(r api.OpinionFiltersRequest) (
// 	entity.OpinionFilters, error) {
// 	var exch *entity.Exchange
// 	if r.SecurityType != nil {
// 		s, err := exchange(*r.Exchange)
// 		if err != nil {
// 			return entity.OpinionFilters{}, err
// 		}
// 		exch = &s
// 	}

// 	var secType *entity.SecurityType
// 	if r.SecurityType != nil {
// 		s, err := securityType(*r.SecurityType)
// 		if err != nil {
// 			return entity.OpinionFilters{}, err
// 		}
// 		secType = &s
// 	}

// 	return entity.OpinionFilters{
// 		ExpertID:     r.ExpertId,
// 		Exchange:     exch,
// 		SecurityID:   r.ExpertId,
// 		SecurityType: secType,
// 	}, nil
// }

// func FromOpinionFiltersReqeustToOpinonFilters(req api.OpinionFiltersRequest) (
// 	entity.OpinionFilters, error) {
// 	var secType *entity.SecurityType
// 	if req.SecurityType != nil {
// 		s, err := securityType(*req.SecurityType)
// 		if err != nil {
// 			return entity.OpinionFilters{}, err
// 		}
// 		secType = &s
// 	}

// 	var exch *entity.Exchange
// 	if req.SecurityType != nil {
// 		e, err := exchange(*req.Exchange)
// 		if err != nil {
// 			return entity.OpinionFilters{}, err
// 		}
// 		exch = &e
// 	}

// 	return entity.OpinionFilters{
// 		ExpertID:     req.ExpertId,
// 		SecurityID:   req.SecurityId,
// 		Exchnage:     exch,
// 		SecurityType: secType,
// 	}, nil
// }
