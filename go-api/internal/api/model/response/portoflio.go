package response

import "github.com/pttrulez/investor-go/internal/model"

type FullPortfolioData struct {
	BondDeals      []*MoexBondDeal  `json:"bondDeals"`
	BondPositions  []*BondPosition  `json:"bondPositions"`
	Cash           int              `json:"cash"`
	Cashouts       []*Cashout       `json:"cashouts"`
	CashoutsSum    int              `json:"cashoutsSum"`
	Compound       bool             `json:"compound"`
	Deals          []*model.Deal    `json:"deals"`
	Deposits       []*Deposit       `json:"deposits"`
	DepositsSum    int              `json:"depositsSum"`
	Name           string           `json:"name"`
	Profitability  int              `json:"profitability"`
	ShareDeals     []*MoexShareDeal `json:"shareDeals"`
	SharePositions []*SharePosition `json:"sharePositions"`
	TotalCost      int              `json:"totalCost"`
}