package domain

type Portfolio struct {
	BondsCost      int
	BondPositions  []Position
	Cash           int
	CashoutsSum    int
	Compound       bool
	Coupons        []Coupon
	CouponsSum     int
	Deals          []Deal
	DepositsSum    int
	Dividends      []Dividend
	DividendsSum   int
	ExpensesSum    int
	ID             int
	Name           string
	Profitability  int
	SharesCost     int
	SharePositions []Position
	TotalCost      int
	Transactions   []Transaction
	UserID         int
}
