package entity

import "time"

type SecurityCommonInfo struct {
	Board     ISSMoexBoard
	Engine    ISSMoexEngine
	ID        int
	LotSize   int
	Market    ISSMoexMarket
	Name      string
	Ticker    string
	ShortName string
}

type Share struct {
	SecurityCommonInfo
	PriceDecimals int
}

type Bond struct {
	SecurityCommonInfo
	CouponPercent   float64
	CouponValue     float64
	CouponFrequency int       // частота выплаты купонов в год
	IssueDate       time.Time // облигации
	FaceValue       int       // номинальная стоимость
	MatDate         time.Time
}

type ISSMoexEngine string

const (
	MoexEngineStock    ISSMoexEngine = "stock"
	MoexEngineCurrency ISSMoexEngine = "currency"
)

type ISSMoexMarket string

const (
	MoexMarketBonds ISSMoexMarket = "bonds" // Рынок облигаций
	// MoexMarketClassica      ISSMoexMarket = "classica"      // Classica
	// MoexMarketCredit        ISSMoexMarket = "credit"        // Рынок кредитов
	// MoexMarketDeposit       ISSMoexMarket = "deposit"       // Депозиты с ЦК
	// MoexMarketForeignShares ISSMoexMarket = "foreignshares" // Иностранные ц.б.
	// MoexMarketForeignNdm    ISSMoexMarket = "foreignndm"    // Иностранные ц.б. РПС
	// MoexMarketGcc           ISSMoexMarket = "gcc"           // РЕПО с ЦК с КСУ
	// MoexMarketIndex         ISSMoexMarket = "index"         // Индексы фондового рынка
	// MoexMarketMoexboard     ISSMoexMarket = "moexboard"     // MOEX Board
	// MoexMarketNdm           ISSMoexMarket = "ndm"           // Режим переговорных сделок
	// MoexMarketNonresCcp     ISSMoexMarket = "nonresccp"     // Рынок РЕПО с ЦК (нерезиденты)
	// MoexMarketNonresNdm     ISSMoexMarket = "nonresndm"     // Режим переговорных сделок (нерезиденты)
	// MoexMarketNonresRepo    ISSMoexMarket = "nonresrepo"    // Рынок РЕПО (нерезиденты)
	// MoexMarketOtc           ISSMoexMarket = "otc"           // ОТС - on the counter?
	// MoexMarketQnv           ISSMoexMarket = "qnv"           // Квал. инвесторы
	// MoexMarketMamc          ISSMoexMarket = "mamc"          // Мультивалютный рынок смешанных активов
	// MoexMarketRepo          ISSMoexMarket = "repo"          // Рынок сделок РЕПО
	// MoexMarketStandard      ISSMoexMarket = "standard"      // Standard
	// MoexMarketSelt          ISSMoexMarket = "selt"          // Валюта: Биржевые сделки с ЦК.
	MoexMarketShares ISSMoexMarket = "shares" // Рынок акций
)

func (e ISSMoexMarket) Validate() bool {
	switch e {
	case MoexMarketBonds:
	// case MoexMarketClassica:
	// case MoexMarketCredit:
	// case MoexMarketDeposit:
	// case MoexMarketForeignShares:
	// case MoexMarketForeignNdm:
	// case MoexMarketGcc:
	// case MoexMarketIndex:
	// case MoexMarketMoexboard:
	// case MoexMarketNdm:
	// case MoexMarketNonresCcp:
	// case MoexMarketNonresNdm:
	// case MoexMarketNonresRepo:
	// case MoexMarketOtc:
	// case MoexMarketQnv:
	// case MoexMarketMamc:
	// case MoexMarketRepo:
	// case MoexMarketStandard:
	// case MoexMarketSelt:
	case MoexMarketShares:
	default:
		return false
	}
	return true
}

type ISSMoexBoard string

const (
	Tqbr ISSMoexBoard = "TQBR" // Т+: Акции и ДР - безадрес
	Cets ISSMoexBoard = "CETS" // Системные сделки - безадрес.
)

type ISSMoexSecurityType string

const (
	CommonShare    ISSMoexSecurityType = "common_share"    // "акция обыкновенная"
	CorporateBond  ISSMoexSecurityType = "corporate_bond"  // "корпоративная облигация"
	ExchangeBond   ISSMoexSecurityType = "exchange_bond"   // "облигация"
	ExchangePpif   ISSMoexSecurityType = "exchange_ppif"   // "биржевой ПИФ"
	Futures        ISSMoexSecurityType = "futures"         // "фьючерс"
	OfzBond        ISSMoexSecurityType = "ofz_bond"        // "ОФЗ"
	PreferredShare ISSMoexSecurityType = "preferred_share" // "акция привелигированная"
	PublicPpif     ISSMoexSecurityType = "public_ppif"     // "публичный ПИФ"
	StockIndex     ISSMoexSecurityType = "stock_index"     // "индекс"
	StockIndexIf   ISSMoexSecurityType = "stock_index_if"  // "iNAV облигаций"
)
