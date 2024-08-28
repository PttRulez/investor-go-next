package converter

import (
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

type WrongEnumValueError struct {
	enumName string
	value    string
}

func (e WrongEnumValueError) Error() string {
	return fmt.Sprintf("wrong %s enum value: %s", e.enumName, e.value)
}

func exchange(e api.Exchange) (entity.Exchange, error) {
	switch e {
	case api.MOEX:
		return entity.EXCHMoex, nil
	default:
		return entity.Exchange(""), WrongEnumValueError{
			enumName: "Exchange",
			value:    string(e),
		}
	}
}

func dealType(d api.DealType) (entity.DealType, error) {
	switch d {
	case api.BUY:
		return entity.DTBuy, nil
	case api.SELL:
		return entity.DTSell, nil
	default:
		return entity.DealType(""), WrongEnumValueError{
			enumName: "DealType",
			value:    string(d),
		}
	}
}

func opinionType(o api.OpinionType) (entity.OpinionType, error) {
	switch o {
	case api.FLAT:
		return entity.Flat, nil
	case api.GENERAL:
		return entity.General, nil
	case api.GROWTH:
		return entity.Growth, nil
	case api.REDUCTION:
		return entity.Reduction, nil
	default:
		return entity.OpinionType(""), WrongEnumValueError{
			enumName: "OpinionType",
			value:    string(o),
		}
	}
}

func securityType(e api.SecurityType) (entity.SecurityType, error) {
	switch e {
	case api.BOND:
		return entity.STBond, nil
	case api.CURRENCY:
		return entity.STCurrency, nil
	case api.FUTURES:
		return entity.STFutures, nil
	case api.INDEX:
		return entity.STIndex, nil
	case api.PIF:
		return entity.STPif, nil
	case api.SHARE:
		return entity.STShare, nil
	default:
		return entity.SecurityType(""), WrongEnumValueError{
			enumName: "SecurityType",
			value:    string(e),
		}
	}
}

func transactionType(t api.TransactionType) (entity.TransactionType, error) {
	switch t {
	case api.CASHOUT:
		return entity.TTCashout, nil
	case api.DEPOSIT:
		return entity.TTDeposit, nil
	default:
		return entity.TransactionType(""), WrongEnumValueError{
			enumName: "TransactionType",
			value:    string(t),
		}
	}
}

// ISS.
func ISSMoexBoardToResponse(d entity.ISSMoexBoard) (api.ISSMoexBoard, error) {
	switch d {
	case entity.Cets:
		return api.CETS, nil
	case entity.Tqbr:
		return api.TQBR, nil
	default:
		return api.ISSMoexBoard(""), WrongEnumValueError{
			enumName: "ISSMoexBoard",
			value:    string(d),
		}
	}
}
func ISSMoexEngineToResponse(d entity.ISSMoexEngine) (api.ISSMoexEngine, error) {
	switch d {
	case entity.MoexEngineCurrency:
		return api.Currency, nil
	case entity.MoexEngineStock:
		return api.Stock, nil
	default:
		return api.ISSMoexEngine(""), WrongEnumValueError{
			enumName: "ISSMoexEngine",
			value:    string(d),
		}
	}
}
func ISSMoexMarketToResponse(d entity.ISSMoexMarket) (api.ISSMoexMarket, error) {
	switch d {
	case entity.MoexMarketBonds:
		return api.Bonds, nil
	case entity.MoexMarketShares:
		return api.Shares, nil
	default:
		return api.ISSMoexMarket(""), WrongEnumValueError{
			enumName: "ISSMoexMarket",
			value:    string(d),
		}
	}
}
