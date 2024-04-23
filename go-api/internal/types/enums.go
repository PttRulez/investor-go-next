package types

type Role string

const (
	Admin    Role = "ADMIN"
	Investor Role = "INVESTOR"
)

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
	ST_Bond  SecurityType = "BOND"
	Currency SecurityType = "CURRENCY"
	Futures  SecurityType = "FUTURES"
	Index    SecurityType = "INDEX"
	Pif      SecurityType = "PIF"
	ST_Share SecurityType = "SHARE"
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
