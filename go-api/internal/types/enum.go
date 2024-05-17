package types

type Exchange string

const (
	EXCH_Moex Exchange = "MOEX"
)

func (e Exchange) Validate() bool {
	switch e {
	case EXCH_Moex:
	default:
		return false
	}
	return true
}

type SecurityType string

const (
	ST_Bond     SecurityType = "BOND"
	ST_Currency SecurityType = "CURRENCY"
	ST_Futures  SecurityType = "FUTURES"
	ST_Index    SecurityType = "INDEX"
	ST_Pif      SecurityType = "PIF"
	ST_Share    SecurityType = "SHARE"
)

func (e SecurityType) Validate() bool {
	switch e {
	case ST_Bond:
	case ST_Share:
	default:
		return false
	}
	return true
}
