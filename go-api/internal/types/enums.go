package types

type Role string

const (
	Admin    Role = "Admin"
	Investor Role = "Investor"
)

type Exchange string

const (
	EXCH_Moex Exchange = "Moex"
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
	ST_Bond  SecurityType = "Bond"
	Currency SecurityType = "Currency"
	Futures  SecurityType = "Futures"
	Index    SecurityType = "Index"
	Pif      SecurityType = "Pif"
	ST_Share SecurityType = "Share"
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
