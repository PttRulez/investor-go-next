package model

type Portfolio struct {
	BondDeals      []*Deal     `json:"bondDeals"`
	BondPositions  []*Position `json:"bondPositions"`
	Cash           int         `json:"cash"`
	Cashouts       []*Cashout  `json:"cashouts"`
	CashoutsSum    int         `json:"cashoutsSum"`
	Compound       bool        `json:"compound"`
	Deals          []*Deal     `json:"deals"`
	Deposits       []*Deposit  `json:"deposits"`
	DepositsSum    int         `json:"depositsSum"`
	Id             int         `json:"id"`
	Name           string      `json:"name"`
	Positions      []*Position `json:"positions"`
	Profitability  int         `json:"profitability"`
	ShareDeals     []*Deal     `json:"sahreDeals"`
	SharePositions []*Position `json:"sharePositions"`
	TotalCost      int         `json:"totalCost"`
	UserId         int         `json:"userId"`
}
