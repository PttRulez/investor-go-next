package model

type MoexBond struct {
	Board     ISSMoexBoard
	Engine    ISSMoexEngine
	Id        int
	Market    ISSMoexMarket
	Name      string
	ShortName string
	Isin      string
}
