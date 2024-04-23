package tmoex

type Engine string

const (
	Stock    = "stock"
	Currency = "currency"
)

type Market string

const (
	Market_Shares = "shares"
	Market_Bonds  = "bonds"
	Market_Index  = "index"
	Market_Selt   = "selt" // Валюта: Биржевые сделки с ЦК
)

type Board string

const (
	Tqbr = "TQBR" // Т+: Акции и ДР - безадрес
	Cets = "CETS" // Системные сделки - безадрес.
)

type SecurityType string

const (
	CommonShare    = "common_share"    // "акция обыкновенная"
	PreferredShare = "preferred_share" // "акция привелигированная"

	CorporateBond = "corporate_bond" // "корпоративная облигация"
	ExchangeBond  = "exchange_bond"  // "облигация"
	OfzBond       = "ofz_bond"       // "ОФЗ"

	ExchangePpif = "exchange_ppif"  // "биржевой ПИФ"
	PublicPpif   = "public_ppif"    // "публичный ПИФ"
	StockIndexIf = "stock_index_if" // "iNAV облигаций"

	Futures = "futures" // "фьючерс"

	StockIndex = "stock_index" // "индекс"
)
