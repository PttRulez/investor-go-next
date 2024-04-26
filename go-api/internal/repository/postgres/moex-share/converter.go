package moexshare

import "github.com/pttrulez/investor-go/internal/model"

func FromMoexShareToDB(share *model.MoexShare) *MoexShare {
	return &MoexShare{
		Id:        share.Id,
		Board:     share.Board,
		Engine:    share.Engine,
		Market:    share.Market,
		Name:      share.Name,
		ShortName: share.ShortName,
		Ticker:    share.Ticker,
	}
}

func FromDBToMoexShare(share *MoexShare) *model.MoexShare {
	return &model.MoexShare{
		Id:        share.Id,
		Board:     share.Board,
		Engine:    share.Engine,
		Market:    share.Market,
		Name:      share.Name,
		ShortName: share.ShortName,
		Ticker:    share.Ticker,
	}
}
