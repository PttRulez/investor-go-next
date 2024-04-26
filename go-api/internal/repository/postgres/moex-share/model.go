package moexshare

import "github.com/pttrulez/investor-go/internal/model"

type MoexShare struct {
	Id        int                 `db:"id"`
	Board     model.ISSMoexBoard  `db:"board"`
	Engine    model.ISSMoexEngine `db:"engine"`
	Market    model.ISSMoexMarket `db:"market"`
	Name      string              `db:"name"`
	ShortName string              `db:"shortname"`
	Ticker    string              `db:"isin"`
}
