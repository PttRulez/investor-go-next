package entity

type Exchange string

const (
	EXCHMoex Exchange = "MOEX"
)

func (e Exchange) Validate() bool {
	switch e {
	case EXCHMoex:
	default:
		return false
	}
	return true
}

type SecurityType string

const (
	STBond     SecurityType = "BOND"
	STCurrency SecurityType = "CURRENCY"
	STFutures  SecurityType = "FUTURES"
	STIndex    SecurityType = "INDEX"
	STPif      SecurityType = "PIF"
	STShare    SecurityType = "SHARE"
)

func (e SecurityType) Validate() bool {
	switch e {
	case STBond:
	case STShare:
	default:
		return false
	}
	return true
}
