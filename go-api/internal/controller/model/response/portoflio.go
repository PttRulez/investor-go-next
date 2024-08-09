package response

type ShortPortfolio struct {
	Id       int    `json:"id"`
	Compound bool   `json:"compound"`
	Name     string `json:"name"`
}

type FullPortfolio struct {
	BondPositions  []Position    `json:"bondPositions"`
	Cash           int           `json:"cash"`
	CashoutsSum    int           `json:"cashoutsSum"`
	Compound       bool          `json:"compound"`
	Deals          []Deal        `json:"deals"`
	DepositsSum    int           `json:"depositsSum"`
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	Profitability  int           `json:"profitability"`
	SharePositions []Position    `json:"sharePositions"`
	TotalCost      int           `json:"totalCost"`
	Transactions   []Transaction `json:"transactions"`
}
