package model

type MoexShare struct {
	Board     ISSMoexBoard
	Engine    ISSMoexEngine
	Id        int
	Market    ISSMoexMarket
	Name      string
	ShortName string
	Ticker    string
}
