package converter

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromMoexBondToMoexBondResponse(b domain.Bond) contracts.MoexBondResponse {
	return contracts.MoexBondResponse{
		Board:           contracts.ISSMoexBoard(b.Board),
		CouponFrequency: b.CouponFrequency,
		CouponPercent:   b.CouponPercent,
		CouponValue:     b.CouponValue,
		Engine:          contracts.ISSMoexEngine(b.Engine),
		FaceValue:       b.FaceValue,
		Id:              b.ID,
		IssueDate:       openapi_types.Date{Time: b.IssueDate},
		LotSize:         b.LotSize,
		Market:          contracts.ISSMoexMarket(b.Market),
		MatDate:         openapi_types.Date{Time: b.MatDate},
		Name:            b.Name,
		Ticker:          b.Ticker,
		ShortName:       b.ShortName,
	}
}

func FromMoexShareToMoexShareResponse(s domain.Share) contracts.MoexShareResponse {
	return contracts.MoexShareResponse{
		Board:     contracts.ISSMoexBoard(s.Board),
		Engine:    contracts.ISSMoexEngine(s.Engine),
		Id:        s.ID,
		LotSize:   s.LotSize,
		Market:    contracts.ISSMoexMarket(s.Market),
		Name:      s.Name,
		Ticker:    s.Ticker,
		ShortName: s.ShortName,
	}
}
