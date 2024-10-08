package converter

import (
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
)

type WrongEnumValueError struct {
	enumName string
	value    string
}

func (e WrongEnumValueError) Error() string {
	return fmt.Sprintf("wrong %s enum value: %s", e.enumName, e.value)
}

func exchange(e contracts.Exchange) (domain.Exchange, error) {
	switch e {
	case contracts.MOEX:
		return domain.EXCHMoex, nil
	default:
		return domain.Exchange(""), WrongEnumValueError{
			enumName: "Exchange",
			value:    string(e),
		}
	}
}

func dealType(d contracts.DealType) (domain.DealType, error) {
	switch d {
	case contracts.BUY:
		return domain.DTBuy, nil
	case contracts.SELL:
		return domain.DTSell, nil
	default:
		return domain.DealType(""), WrongEnumValueError{
			enumName: "DealType",
			value:    string(d),
		}
	}
}

func opinionType(o contracts.OpinionType) (domain.OpinionType, error) {
	switch o {
	case contracts.FLAT:
		return domain.Flat, nil
	case contracts.GENERAL:
		return domain.General, nil
	case contracts.GROWTH:
		return domain.Growth, nil
	case contracts.REDUCTION:
		return domain.Reduction, nil
	default:
		return domain.OpinionType(""), WrongEnumValueError{
			enumName: "OpinionType",
			value:    string(o),
		}
	}
}

func securityType(e contracts.SecurityType) (domain.SecurityType, error) {
	switch e {
	case contracts.BOND:
		return domain.STBond, nil
	case contracts.CURRENCY:
		return domain.STCurrency, nil
	case contracts.FUTURES:
		return domain.STFutures, nil
	case contracts.INDEX:
		return domain.STIndex, nil
	case contracts.PIF:
		return domain.STPif, nil
	case contracts.SHARE:
		return domain.STShare, nil
	default:
		return domain.SecurityType(""), WrongEnumValueError{
			enumName: "SecurityType",
			value:    string(e),
		}
	}
}

func transactionType(t contracts.TransactionType) (domain.TransactionType, error) {
	switch t {
	case contracts.CASHOUT:
		return domain.TTCashout, nil
	case contracts.DEPOSIT:
		return domain.TTDeposit, nil
	default:
		return domain.TransactionType(""), WrongEnumValueError{
			enumName: "TransactionType",
			value:    string(t),
		}
	}
}

// ISS.
func ISSMoexBoardToResponse(d domain.ISSMoexBoard) (contracts.ISSMoexBoard, error) {
	switch d {
	case domain.Cets:
		return contracts.CETS, nil
	case domain.Tqbr:
		return contracts.TQBR, nil
	default:
		return contracts.ISSMoexBoard(""), WrongEnumValueError{
			enumName: "ISSMoexBoard",
			value:    string(d),
		}
	}
}
func ISSMoexEngineToResponse(d domain.ISSMoexEngine) (contracts.ISSMoexEngine, error) {
	switch d {
	case domain.MoexEngineCurrency:
		return contracts.Currency, nil
	case domain.MoexEngineStock:
		return contracts.Stock, nil
	default:
		return contracts.ISSMoexEngine(""), WrongEnumValueError{
			enumName: "ISSMoexEngine",
			value:    string(d),
		}
	}
}
func ISSMoexMarketToResponse(d domain.ISSMoexMarket) (contracts.ISSMoexMarket, error) {
	switch d {
	case domain.MoexMarketBonds:
		return contracts.Bonds, nil
	case domain.MoexMarketShares:
		return contracts.Shares, nil
	default:
		return contracts.ISSMoexMarket(""), WrongEnumValueError{
			enumName: "ISSMoexMarket",
			value:    string(d),
		}
	}
}
