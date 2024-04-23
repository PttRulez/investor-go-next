package types

import "net/http"

type AuthController interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

type CashoutController interface {
	CreateNewCashout(w http.ResponseWriter, r *http.Request)
	DeleteCashout(w http.ResponseWriter, r *http.Request)
}
type MoexShareController interface {
	GetInfoByTicker(w http.ResponseWriter, r *http.Request)
}
type MoexBondController interface {
	GetInfoByISIN(w http.ResponseWriter, r *http.Request)
}
type MoexBondDealController interface {
	CreateNewDeal(w http.ResponseWriter, r *http.Request)
	DeleteDeal(w http.ResponseWriter, r *http.Request)
}
type MoexShareDealController interface {
	CreateNewDeal(w http.ResponseWriter, r *http.Request)
	DeleteDeal(w http.ResponseWriter, r *http.Request)
}

type DepositController interface {
	CreateNewDeposit(w http.ResponseWriter, r *http.Request)
	DeleteDeposit(w http.ResponseWriter, r *http.Request)
}

type ExpertController interface {
	CreateNewExpert(w http.ResponseWriter, r *http.Request)
	GetExpertsList(w http.ResponseWriter, r *http.Request)
}

type PortfolioController interface {
	CreateNewPortfolio(w http.ResponseWriter, r *http.Request)
	GetPortfolioById(w http.ResponseWriter, r *http.Request)
	DeletePortfolio(w http.ResponseWriter, r *http.Request)
	GetListOfPortfolios(w http.ResponseWriter, r *http.Request)
	UpdatePortfolio(w http.ResponseWriter, r *http.Request)
}
