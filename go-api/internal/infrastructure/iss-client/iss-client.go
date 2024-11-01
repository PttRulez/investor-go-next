package issclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
)

type ISSSecInfo struct {
	Board     domain.ISSMoexBoard
	Currency  string
	Engine    domain.ISSMoexEngine
	Market    domain.ISSMoexMarket
	Name      string
	ShortName string
	Ticker    string

	// Только облиги:
	//
	// частота выплаты купонов в год
	CouponFrequency int

	CouponPercent float64
	CouponValue   float64

	// номинальная стоимость облигации
	LotPrice int

	IssueDate time.Time

	// дата погашения облиги
	MaturityDate time.Time
}

type ISSFullSecurityInfo struct {
	LotPrice      int
	PriceDecimals int
}

const (
	couponFrequency = "COUPONFREQUENCY"
	couponPercent   = "COUPONPERCENT"
	couponValue     = "COUPONVALUE"
	currency        = "FACEUNIT"
	issueDate       = "ISSUEDATE"
	lotAmount       = "FACEVALUE"
	maturityDate    = "MATDATE"
	name            = "NAME"
	shortName       = "SHORTNAME"

	usdRateName = "CBRF_USD_LAST"
	eurRateName = "CBRF_EUR_LAST"
)

func (api *IssClient) GetSecurityInfoByTicker(ctx context.Context, ticker string) (ISSSecInfo, error) {
	const op = "issclient.GetSecurityInfoByTicker"

	uri := fmt.Sprintf("%s/securities/%s.json", api.baseURL, ticker)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return ISSSecInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	// фильтруем только то что нам нужно
	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine,is_primary")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return ISSSecInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var data MoexAPIResponseSecurityInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ISSSecInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	result := ISSSecInfo{}
	for _, boardData := range data.Boards.Data {
		isPrimaryFloat, ok := boardData[3].(float64)
		if !ok {
			return ISSSecInfo{}, fmt.Errorf("%s failed to typecast boardData.is_primary = %v", op, boardData[3])
		}
		isPrimary := int(isPrimaryFloat)
		if isPrimary != 1 {
			continue
		}

		board, ok := boardData[0].(string)
		if !ok {
			return ISSSecInfo{}, fmt.Errorf("%s  failed to typecast board = %v", op, boardData[0])
		}

		market, ok := boardData[1].(string)
		if !ok {
			return ISSSecInfo{}, fmt.Errorf("%s failed to typecast market = %v", op, boardData[1])
		}

		engine, ok := boardData[2].(string)
		if !ok {
			return ISSSecInfo{}, fmt.Errorf("%s failed to typecast engine = %v", op, boardData[2])
		}

		result.Board = domain.ISSMoexBoard(board)
		result.Market = domain.ISSMoexMarket(market)
		result.Engine = domain.ISSMoexEngine(engine)
		break
	}

	for _, item := range data.Description.Data {
		err := chooseSecProperty(&result, item[0], item[1])
		if err != nil {
			return ISSSecInfo{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return result, nil
}

func chooseSecProperty(sec *ISSSecInfo, propName, value string) error {
	switch propName {
	case currency:
		if value == "SUR" {
			value = "RUB"
		}
		sec.Currency = value
	case name:
		fmt.Println("NAME", value)
		sec.Name = value
	case shortName:
		sec.ShortName = value
	}

	if sec.Market == domain.MoexMarketBonds {
		switch propName {
		case couponFrequency:
			freq, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", couponFrequency, err)
			}
			sec.CouponFrequency = freq
		case couponPercent:
			percent, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", couponPercent, err)
			}
			sec.CouponPercent = percent
		case couponValue:
			percent, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", couponValue, err)
			}
			sec.CouponValue = percent

		case issueDate:
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", issueDate, err)
			}
			sec.IssueDate = t
		case lotAmount:
			faceValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", lotAmount, err)
			}
			sec.LotPrice = faceValue
		case maturityDate:
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return fmt.Errorf("%s parsing: %w", maturityDate, err)
			}
			sec.MaturityDate = t
		}
	}

	return nil
}

type MoexFullInfo struct {
	Securities struct {
		Data [][]any `json:"data"`
	} `json:"securities"`
}

func (api *IssClient) GetSecurityFullInfo(ctx context.Context, engine domain.ISSMoexEngine,
	market domain.ISSMoexMarket, board domain.ISSMoexBoard, ticker string) (ISSFullSecurityInfo, error) {
	const op = "issclient.GetSecurityFullInfo"

	uri := fmt.Sprintf("%s/engines/%s/markets/%s/boards/%s/securities/%s.json", api.baseURL,
		engine, market, board, ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return ISSFullSecurityInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	// фильтруем только то что нам нужно
	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("securities.columns", "LOTSIZE,PREVPRICE")
	params.Add("marketdata.columns", "off")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return ISSFullSecurityInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var data MoexFullInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ISSFullSecurityInfo{}, fmt.Errorf("%s (failed to Unmarshal) : %w", op, err)
	}

	var result ISSFullSecurityInfo
	lotSizeFloat, ok := data.Securities.Data[0][0].(float64)
	if !ok {
		return ISSFullSecurityInfo{}, fmt.Errorf("%s (failed to TypeCast lotSize (%v) from Data) : %w",
			op, data.Securities.Data[0][0], err)
	}
	prevPrice, ok := data.Securities.Data[0][1].(float64)
	if !ok {
		return ISSFullSecurityInfo{}, fmt.Errorf("%s (failed to TypeCast prevPrice (%v) from Data) : %w",
			op, data.Securities.Data[0][1], err)
	}

	lotSize := int(lotSizeFloat)

	result.LotPrice = lotSize
	result.PriceDecimals = utils.SignsAfterDot(prevPrice)

	return result, nil
}

func (api *IssClient) GetStocksCurrentPrices(ctx context.Context, market domain.ISSMoexMarket,
	tickerInfos map[string]domain.ISSMoexBoard) (map[string]float64, error) {
	const op = "issclient.GetStocksCurrentPrices"

	uri := fmt.Sprintf("%s/engines/stock/markets/%s/securities.json", api.baseURL, market)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tickers := make([]string, 0, len(tickerInfos))
	for ticker := range tickerInfos {
		tickers = append(tickers, ticker)
	}

	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("securities", strings.Join(tickers, ","))
	params.Add("securities.columns", "SECID,BOARDID,PREVPRICE")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var data MoexAPIResponseCurrentPrices
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var m = make(map[string]float64)
	for _, i := range data.Securities.Data {
		ticker, ok := i[0].(string)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast ticker from issreponse", op)
		}

		board, ok := i[1].(string)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast %v board from issreponse", i[1], op)
		}

		price, ok := i[2].(float64)
		if !ok {
			// цена от iss иногда приходит null для некоторых бордов
			if i[2] == nil {
				continue
			}
			return nil, fmt.Errorf("%s: failed to cast price from issreponse", op)
		}
		v, ok := tickerInfos[ticker]
		if ok && string(v) == board {
			m[ticker] = price
		}
	}

	return m, nil
}

func (api *IssClient) GetCurrencyRates(ctx context.Context) (map[string]float64, error) {
	const op = "issclient.GetCurrencyRates"

	uri := fmt.Sprintf("%s/statistics/engines/currency/markets/selt/rates.json",
		api.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s NewRequestWithContext: %w", op, err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s api.client.Do: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var data MoexAPIResponseCurrency
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := map[string]float64{"RUB": 1}
	for i, name := range data.Cbrf.RateNames {
		if name == usdRateName {
			rate, ok := data.Cbrf.Data[0][i].(float64)
			if !ok {
				return nil, fmt.Errorf("%s: failed to cast price from issreponse", op)
			}
			result["USD"] = rate
		} else if name == eurRateName {
			rate, ok := data.Cbrf.Data[0][i].(float64)
			if !ok {
				return nil, fmt.Errorf("%s: failed to cast price from issreponse", op)
			}
			result["EUR"] = rate
		}
	}

	return result, nil
}

type IssClient struct {
	baseURL string
	client  *http.Client
}

func NewISSClient() *IssClient {
	return &IssClient{
		baseURL: "https://iss.moex.com/iss",
		client:  &http.Client{},
	}
}
