package response

type ShortPortfolio struct {
	Id       int    `json:"id"`
	Compound bool   `json:"compound"`
	Name     string `json:"name"`
}

type FullPortfolio struct {
	BondDeals      []Deal        `json:"bondDeals"`
	BondPositions  []Position    `json:"bondPositions"`
	Cash           int           `json:"cash"`
	Cashouts       []Cashout     `json:"cashouts"`
	CashoutsSum    int           `json:"cashoutsSum"`
	Compound       bool          `json:"compound"`
	Deals          []interface{} `json:"deals"`
	Deposits       []Deposit     `json:"deposits"`
	DepositsSum    int           `json:"depositsSum"`
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	Profitability  int           `json:"profitability"`
	ShareDeals     []Deal        `json:"shareDeals"`
	SharePositions []Position    `json:"sharePositions"`
	TotalCost      int           `json:"totalCost"`
}
