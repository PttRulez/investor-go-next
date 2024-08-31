// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for DealType.
const (
	BUY  DealType = "BUY"
	SELL DealType = "SELL"
)

// Defines values for Exchange.
const (
	MOEX Exchange = "MOEX"
)

// Defines values for ISSMoexBoard.
const (
	CETS ISSMoexBoard = "CETS"
	TQBR ISSMoexBoard = "TQBR"
)

// Defines values for ISSMoexEngine.
const (
	Currency ISSMoexEngine = "currency"
	Stock    ISSMoexEngine = "stock"
)

// Defines values for ISSMoexMarket.
const (
	Bonds  ISSMoexMarket = "bonds"
	Shares ISSMoexMarket = "shares"
)

// Defines values for OpinionType.
const (
	FLAT      OpinionType = "FLAT"
	GENERAL   OpinionType = "GENERAL"
	GROWTH    OpinionType = "GROWTH"
	REDUCTION OpinionType = "REDUCTION"
)

// Defines values for Role.
const (
	ADMIN    Role = "ADMIN"
	INVESTOR Role = "INVESTOR"
)

// Defines values for SecurityType.
const (
	BOND     SecurityType = "BOND"
	CURRENCY SecurityType = "CURRENCY"
	FUTURES  SecurityType = "FUTURES"
	INDEX    SecurityType = "INDEX"
	PIF      SecurityType = "PIF"
	SHARE    SecurityType = "SHARE"
)

// Defines values for TransactionType.
const (
	CASHOUT TransactionType = "CASHOUT"
	DEPOSIT TransactionType = "DEPOSIT"
)

// CreateDealRequest defines model for CreateDealRequest.
type CreateDealRequest struct {
	Amount       int          `json:"amount"`
	Comission    float64      `json:"comission"`
	Date         time.Time    `json:"date"`
	Exchange     Exchange     `json:"exchange" validate:"required,is-exchange"`
	PortfolioId  int          `json:"portfolioId"`
	Price        float64      `json:"price"`
	SecurityType SecurityType `json:"securityType" validate:"required,securityType"`
	Ticker       string       `json:"ticker"`
	Type         DealType     `json:"type"`
}

// CreateExpertRequest defines model for CreateExpertRequest.
type CreateExpertRequest struct {
	AvatarUrl *string `json:"avatarUrl,omitempty"`
	Name      string  `json:"name"`
}

// CreateOpinionRequest defines model for CreateOpinionRequest.
type CreateOpinionRequest struct {
	Date         openapi_types.Date `json:"date" validate:"required"`
	Exchange     Exchange           `json:"exchange" validate:"required,is-exchange"`
	ExpertId     int                `json:"expertId" validate:"required"`
	SecurityId   int                `json:"securityId" validate:"required"`
	SecurityType SecurityType       `json:"securityType" validate:"required,securityType"`
	SourceLink   *string            `json:"sourceLink,omitempty"`
	TargetPrice  *float64           `json:"targetPrice,omitempty"`
	Text         string             `json:"text" validate:"required"`
	Ticker       string             `json:"ticker" validate:"required"`
	Type         OpinionType        `json:"type" validate:"required,opinionType"`
}

// CreatePortfolioRequest defines model for CreatePortfolioRequest.
type CreatePortfolioRequest struct {
	Compound bool   `json:"compound"`
	Name     string `json:"name"`
}

// CreateTransactionRequest defines model for CreateTransactionRequest.
type CreateTransactionRequest struct {
	Amount      int             `json:"amount"`
	Date        time.Time       `json:"date"`
	PortfolioId int             `json:"portfolioId"`
	Type        TransactionType `json:"type"`
}

// DealResponse defines model for DealResponse.
type DealResponse struct {
	Amount       int          `json:"amount"`
	Comission    float64      `json:"comission"`
	Date         time.Time    `json:"date"`
	Exchange     Exchange     `json:"exchange" validate:"required,is-exchange"`
	Id           *int         `json:"id,omitempty"`
	PortfolioId  int          `json:"portfolioId"`
	Price        float64      `json:"price"`
	SecurityType SecurityType `json:"securityType" validate:"required,securityType"`
	Ticker       string       `json:"ticker"`
	Type         DealType     `json:"type"`
}

// DealType defines model for DealType.
type DealType string

// Exchange defines model for Exchange.
type Exchange string

// ExpertListResponse defines model for ExpertListResponse.
type ExpertListResponse = []ExpertResponse

// ExpertResponse defines model for ExpertResponse.
type ExpertResponse struct {
	AvatarUrl *string `json:"avatarUrl,omitempty"`
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	UserId    *int    `json:"userId,omitempty"`
}

// FullPortfolioResponse defines model for FullPortfolioResponse.
type FullPortfolioResponse struct {
	BondPositions  []PositionResponse    `json:"bondPositions"`
	Cash           int                   `json:"cash"`
	CashoutsSum    int                   `json:"cashoutsSum"`
	Compound       bool                  `json:"compound"`
	Deals          []DealResponse        `json:"deals"`
	DepositsSum    int                   `json:"depositsSum"`
	Id             int                   `json:"id"`
	Name           string                `json:"name"`
	Profitability  int                   `json:"profitability"`
	SharePositions []PositionResponse    `json:"sharePositions"`
	TotalCost      int                   `json:"totalCost"`
	Transactions   []TransactionResponse `json:"transactions"`
}

// ISSMoexBoard defines model for ISSMoexBoard.
type ISSMoexBoard string

// ISSMoexEngine defines model for ISSMoexEngine.
type ISSMoexEngine string

// ISSMoexMarket defines model for ISSMoexMarket.
type ISSMoexMarket string

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token *string `json:"token,omitempty"`
}

// MoexBondResponse defines model for MoexBondResponse.
type MoexBondResponse struct {
	Board ISSMoexBoard `json:"board"`

	// CouponFrequency частота выплаты купонов в год
	CouponFrequency int `json:"couponFrequency"`

	// CouponPercent купон в процентах
	CouponPercent float64 `json:"couponPercent"`

	// CouponValue купон в деньгах
	CouponValue float64       `json:"couponValue"`
	Engine      ISSMoexEngine `json:"engine"`

	// FaceValue номинальная стоимость
	FaceValue int                `json:"faceValue"`
	Id        int                `json:"id"`
	IssueDate openapi_types.Date `json:"issueDate"`
	LotSize   int                `json:"lotSize"`
	Market    ISSMoexMarket      `json:"market"`
	MatDate   openapi_types.Date `json:"matDate"`
	Name      string             `json:"name"`
	ShortName string             `json:"shortName"`
	Ticker    string             `json:"ticker"`
}

// MoexSecurityResponse defines model for MoexSecurityResponse.
type MoexSecurityResponse struct {
	Board     ISSMoexBoard  `json:"board"`
	Engine    ISSMoexEngine `json:"engine"`
	Id        int           `json:"id"`
	LotSize   int           `json:"lotSize"`
	Market    ISSMoexMarket `json:"market"`
	Name      string        `json:"name"`
	ShortName string        `json:"shortName"`
	Ticker    string        `json:"ticker"`
}

// MoexShareResponse defines model for MoexShareResponse.
type MoexShareResponse = MoexSecurityResponse

// OpinionResponse defines model for OpinionResponse.
type OpinionResponse struct {
	Date         openapi_types.Date `json:"date"`
	Exchange     Exchange           `json:"exchange" validate:"required,is-exchange"`
	Expert       ExpertResponse     `json:"expert"`
	ExpertId     int                `json:"expertId"`
	Id           int                `json:"id"`
	SecurityId   int                `json:"securityId"`
	SecurityType SecurityType       `json:"securityType" validate:"required,securityType"`
	SourceLink   *string            `json:"sourceLink,omitempty"`
	TargetPrice  *float64           `json:"targetPrice,omitempty"`
	Text         string             `json:"text"`
	Ticker       string             `json:"ticker"`
	Type         OpinionType        `json:"type" validate:"required,opinionType"`
}

// OpinionType defines model for OpinionType.
type OpinionType string

// PortfolioListResponse defines model for PortfolioListResponse.
type PortfolioListResponse = []PortfolioResponse

// PortfolioResponse defines model for PortfolioResponse.
type PortfolioResponse struct {
	Compound bool   `json:"compound"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UserId   *int   `json:"userId,omitempty"`
}

// PositionResponse defines model for PositionResponse.
type PositionResponse struct {
	Amount        int               `json:"amount"`
	AveragePrice  float64           `json:"averagePrice"`
	Comment       *string           `json:"comment,omitempty"`
	CurrentCost   int               `json:"currentCost"`
	CurrentPrice  float64           `json:"currentPrice"`
	Id            int               `json:"id"`
	OpinionIds    []int             `json:"opinionIds"`
	Opinions      []OpinionResponse `json:"opinions"`
	PortfolioName string            `json:"portfolioName"`
	SecurityType  SecurityType      `json:"securityType" validate:"required,securityType"`
	ShortName     string            `json:"shortName"`
	TargetPrice   *float64          `json:"targetPrice,omitempty"`
	Ticker        string            `json:"ticker"`
}

// RegisterUserRequest defines model for RegisterUserRequest.
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     *Role  `json:"role,omitempty"`
}

// Role defines model for Role.
type Role string

// SecurityType defines model for SecurityType.
type SecurityType string

// TransactionResponse defines model for TransactionResponse.
type TransactionResponse struct {
	Amount      int             `json:"amount"`
	Date        time.Time       `json:"date"`
	Id          int             `json:"id"`
	PortfolioId int             `json:"portfolioId"`
	Type        TransactionType `json:"type"`
}

// TransactionType defines model for TransactionType.
type TransactionType string

// UpdatePortfolioRequest defines model for UpdatePortfolioRequest.
type UpdatePortfolioRequest struct {
	Compound *bool   `json:"compound,omitempty"`
	Id       int     `json:"id"`
	Name     *string `json:"name,omitempty"`
}

// GetOpinionListParams defines parameters for GetOpinionList.
type GetOpinionListParams struct {
	// ExpertId id эксперта
	ExpertId *int `form:"expertId,omitempty" json:"expertId,omitempty"`

	// SecurityId id бумаги
	SecurityId *int `form:"securityId,omitempty" json:"securityId,omitempty"`

	// Exchange Биржа
	Exchange *string `form:"exchange,omitempty" json:"exchange,omitempty"`

	// SecurityType тип бумаги SHARE BOND ...
	SecurityType *string `form:"securityType,omitempty" json:"securityType,omitempty"`
}

// PostDealJSONRequestBody defines body for PostDeal for application/json ContentType.
type PostDealJSONRequestBody = CreateDealRequest

// PostExpertJSONRequestBody defines body for PostExpert for application/json ContentType.
type PostExpertJSONRequestBody = CreateExpertRequest

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody = LoginRequest

// PostOpinionJSONRequestBody defines body for PostOpinion for application/json ContentType.
type PostOpinionJSONRequestBody = CreateOpinionRequest

// PostPortfolioJSONRequestBody defines body for PostPortfolio for application/json ContentType.
type PostPortfolioJSONRequestBody = CreatePortfolioRequest

// PutPortfolioJSONRequestBody defines body for PutPortfolio for application/json ContentType.
type PutPortfolioJSONRequestBody = UpdatePortfolioRequest

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody = RegisterUserRequest

// PostTransactionJSONRequestBody defines body for PostTransaction for application/json ContentType.
type PostTransactionJSONRequestBody = CreateTransactionRequest
