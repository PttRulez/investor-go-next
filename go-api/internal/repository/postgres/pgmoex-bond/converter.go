package pgmoexbond

import "github.com/pttrulez/investor-go/internal/model"

func FromMoexBondToDB(bond *model.MoexBond) *MoexBond {
	return &MoexBond{
		Id:        bond.Id,
		Board:     bond.Board,
		Engine:    bond.Engine,
		Market:    bond.Market,
		Name:      bond.Name,
		ShortName: bond.ShortName,
		Isin:      bond.Isin,
	}
}

func FromDBToMoexBond(bond *MoexBond) *model.MoexBond {
	return &model.MoexBond{
		Id:        bond.Id,
		Board:     bond.Board,
		Engine:    bond.Engine,
		Market:    bond.Market,
		Name:      bond.Name,
		ShortName: bond.ShortName,
		Isin:      bond.Isin,
	}
}
