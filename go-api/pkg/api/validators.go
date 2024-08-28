package api

func (e ISSMoexMarket) Validate() bool {
	switch e {
	case Bonds:
	case Shares:
	default:
		return false
	}
	return true
}

func (e DealType) Validate() bool {
	switch e {
	case BUY:
	case SELL:
	default:
		return false
	}
	return true
}

func (e Exchange) Validate() bool {
	switch e {
	case MOEX:
	default:
		return false
	}
	return true
}

func (e OpinionType) Validate() bool {
	switch e {
	case FLAT:
	case GENERAL:
	case GROWTH:
	case REDUCTION:
	default:
		return false
	}
	return true
}

func (e TransactionType) Validate() bool {
	switch e {
	case DEPOSIT:
	case CASHOUT:
	default:
		return false
	}
	return true
}

func (e SecurityType) Validate() bool {
	switch e {
	case BOND:
	case CURRENCY:
	case FUTURES:
	case INDEX:
	case PIF:
	case SHARE:
	default:
		return false
	}
	return true
}
