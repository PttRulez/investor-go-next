package entity

type Portfolio struct {
	BondPositions  []Position
	Cash           int
	CashoutsSum    int
	Compound       bool
	Deals          []Deal
	DepositsSum    int
	ID             int
	Name           string
	Profitability  int
	SharePositions []Position
	TotalCost      int
	Transactions   []Transaction
	UserID         int
}
