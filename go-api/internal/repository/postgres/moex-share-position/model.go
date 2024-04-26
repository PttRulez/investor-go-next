package moexshareposition

type MoexSharePosition struct {
	Id           int     `db:"id"`
	Amount       int     `db:"amount"`
	AveragePrice float64 `db:"average_price"`
	Comment      string  `db:"comment"`
	PortfolioId  int     `db:"portfolio_id"`
	SecurityId   int     `db:"security_id"`
	TargetPrice  float64 `db:"target_price"`
	Ticker       string  `db:"isin"`
	ShortName    string  `db:"shortname"`
}
