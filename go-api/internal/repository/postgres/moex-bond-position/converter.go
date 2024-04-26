package moexbondposition

import "github.com/pttrulez/investor-go/internal/model"

func FromDBtoPosition(db *MoexBondPosition) *model.Position {
	return &model.Position{
		Amount:       db.Amount,
		AveragePrice: db.AveragePrice,
		Comment:      db.Comment,
		Exchange:     model.EXCH_Moex,
		PortfolioId:  db.PortfolioId,
		Secid:        db.Isin,
		SecurityId:   db.SecurityId,
		SecurityType: model.ST_Bond,
		ShortName:    db.ShortName,
		TargetPrice:  db.TargetPrice,
	}
}

func FromPositionToDB(position *model.Position) *MoexBondPosition {
	return &MoexBondPosition{
		Amount:       position.Amount,
		AveragePrice: position.AveragePrice,
		Comment:      position.Comment,
		Isin:         position.Secid,
		PortfolioId:  position.PortfolioId,
		SecurityId:   position.SecurityId,
		ShortName:    position.ShortName,
		TargetPrice:  position.TargetPrice,
	}
}
