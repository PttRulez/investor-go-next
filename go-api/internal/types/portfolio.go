package types

type Portfolio struct {
	Compound bool   `json:"compound" db:"compound"`
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"  validate:"required"`
	UserId   int    `json:"-" db:"user_id"`
}

type PortfolioUpdate struct {
	Compound *bool   `json:"compound,omitempty" db:"compound" validate:"bool"`
	Name     *string `json:"name,omitempty" db:"name"`
	Id       int     `json:"id" db:"id"  validate:"required"`
}

type FullPortfolioData struct {
	Portfolio
	BondPositions  []*BondPosition  `json:"bondPositions"`
	Cash           int              `json:"cash"`
	Cashouts       []*Cashout       `json:"cashouts"`
	CashoutsSum    int              `json:"cashoutsSum"`
	ShareDeals     []*Deal          `json:"shareDeals"`
	BondDeals      []*Deal          `json:"bondDeals"`
	Deposits       []*Deposit       `json:"deposits"`
	DepositsSum    int              `json:"depositsSum"`
	SharePositions []*SharePosition `json:"sharePositions"`
	Porfitability  float64          `json:"porfitability"`
	TotalCost      int              `json:"totalCost"`
}
