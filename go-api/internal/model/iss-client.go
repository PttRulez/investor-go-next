package model

type MoexApiResponseSecurityInfo struct {
	Description struct {
		Data [][]string `json:"data"`
	}
	Boards struct {
		Data [][4]any `json:"data"`
	}
}

type MoexApiResponseCurrentPrices struct {
	Securities struct {
		Data [][3]any
	}
}

type ISSMoexEngine string

const (
	ME_Stock    ISSMoexEngine = "stock"
	ME_Currency ISSMoexEngine = "currency"
)

type ISSMoexMarket string

const (
	Market_Bonds         ISSMoexMarket = "bonds"         // Рынок облигаций
	Market_Classica      ISSMoexMarket = "classica"      // Classica
	Market_Credit        ISSMoexMarket = "credit"        // Рынок кредитов
	Market_Deposit       ISSMoexMarket = "deposit"       // Депозиты с ЦК
	Market_ForeignShares ISSMoexMarket = "foreignshares" // Иностранные ц.б.
	Market_ForeignNdm    ISSMoexMarket = "foreignndm"    // Иностранные ц.б. РПС
	Market_Gcc           ISSMoexMarket = "gcc"           // РЕПО с ЦК с КСУ
	Market_Index         ISSMoexMarket = "index"         // Индексы фондового рынка
	Market_Moexboard     ISSMoexMarket = "moexboard"     // MOEX Board
	Market_Ndm           ISSMoexMarket = "ndm"           // Режим переговорных сделок
	Market_NonresCcp     ISSMoexMarket = "nonresccp"     // Рынок РЕПО с ЦК (нерезиденты)
	Market_NonresNdm     ISSMoexMarket = "nonresndm"     // Режим переговорных сделок (нерезиденты)
	Market_NonresRepo    ISSMoexMarket = "nonresrepo"    // Рынок РЕПО (нерезиденты)
	Market_Otc           ISSMoexMarket = "otc"           // ОТС - on the counter?
	Market_Qnv           ISSMoexMarket = "qnv"           // Квал. инвесторы
	Market_Mamc          ISSMoexMarket = "mamc"          //Мультивалютный рынок смешанных активов
	Market_Repo          ISSMoexMarket = "repo"          // Рынок сделок РЕПО
	Market_Standard      ISSMoexMarket = "standard"      // Standard
	Market_Selt          ISSMoexMarket = "selt"          // Валюта: Биржевые сделки с ЦК
	Market_Shares        ISSMoexMarket = "shares"        // Рынок акций
)

func (e ISSMoexMarket) Validate() bool {
	switch e {
	case Market_Bonds:
	case Market_Classica:
	case Market_Credit:
	case Market_Deposit:
	case Market_ForeignShares:
	case Market_ForeignNdm:
	case Market_Gcc:
	case Market_Index:
	case Market_Moexboard:
	case Market_Ndm:
	case Market_NonresCcp:
	case Market_NonresNdm:
	case Market_NonresRepo:
	case Market_Otc:
	case Market_Qnv:
	case Market_Mamc:
	case Market_Repo:
	case Market_Standard:
	case Market_Selt:
	case Market_Shares:
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
