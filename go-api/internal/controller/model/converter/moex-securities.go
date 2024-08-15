package converter

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromMoexBondToMoexBondResponse(b entity.Bond) api.MoexBondResponse {
	return api.MoexBondResponse{
		Board:           api.ISSMoexBoard(b.Board),
		CouponFrequency: b.CouponFrequency,
		CouponPercent:   b.CouponPercent,
		CouponValue:     b.CouponValue,
		Engine:          api.ISSMoexEngine(b.Engine),
		FaceValue:       b.FaceValue,
		Id:              b.ID,
		IssueDate:       openapi_types.Date{Time: b.IssueDate},
		Market:          api.ISSMoexMarket(b.Market),
		MatDate:         openapi_types.Date{Time: b.MatDate},
		Name:            b.Name,
		Secid:           b.Secid,
		ShortName:       b.ShortName,
	}
}

func FromMoexShareToMoexShareResponse(s entity.Share) api.MoexShareResponse {
	return api.MoexShareResponse{
		Board:     api.ISSMoexBoard(s.Board),
		Engine:    api.ISSMoexEngine(s.Engine),
		Id:        s.ID,
		Market:    api.ISSMoexMarket(s.Market),
		Name:      s.Name,
		Secid:     s.Secid,
		ShortName: s.ShortName,
	}
}
