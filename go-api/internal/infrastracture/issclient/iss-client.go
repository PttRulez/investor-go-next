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

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
)

type ISSSecInfo struct {
	Board     entity.ISSMoexBoard
	Engine    entity.ISSMoexEngine
	Market    entity.ISSMoexMarket
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
	FaceValue int

	IssueDate time.Time

	// дата погашения облиги
	MatDate time.Time
}

type ISSFullSecurityInfo struct {
	LotSize       int
	PriceDecimals int
}

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

		result.Board = entity.ISSMoexBoard(board)
		result.Market = entity.ISSMoexMarket(market)
		result.Engine = entity.ISSMoexEngine(engine)
		break
	}

	// У облигаций больше свойств
	if result.Market == entity.MoexMarketBonds {
		for _, item := range data.Description.Data {
			switch item[0] {
			case "NAME":
				result.Name = item[1]
			case "SHORTNAME":
				result.ShortName = item[1]
			case "COUPONFREQUENCY":
				freq, err := strconv.Atoi(item[1])
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing COUPONFREQUENCY: %w", op, err)
				}
				result.CouponFrequency = freq
			case "COUPONPERCENT":
				percent, err := strconv.ParseFloat(item[1], 32)
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing COUPONPERCENT: %w", op, err)
				}
				result.CouponPercent = percent
			case "COUPONVALUE":
				percent, err := strconv.ParseFloat(item[1], 32)
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing COUPONVALUE: %w", op, err)
				}
				result.CouponValue = percent
			case "ISSUEDATE":
				t, err := time.Parse("2006-01-02", item[1])
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing ISSUEDATE: %w", op, err)
				}
				result.IssueDate = t
			case "MATDATE":
				t, err := time.Parse("2006-01-02", item[1])
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing MATDATE: %w", op, err)
				}
				result.MatDate = t
			case "FACEVALUE":
				faceValue, err := strconv.Atoi(item[1])
				if err != nil {
					return ISSSecInfo{}, fmt.Errorf("%s parsing FACEVALUE: %w", op, err)
				}
				result.FaceValue = faceValue
			}
		}
	} else {
		for _, item := range data.Description.Data {
			switch item[0] {
			case "NAME":
				result.Name = item[1]
			case "SHORTNAME":
				result.ShortName = item[1]
			}
		}
	}

	return result, nil
}

type MoexFullInfo struct {
	Securities struct {
		Data [][]any `json:"data"`
	} `json:"securities"`
}

func (api *IssClient) GetSecurityFullInfo(ctx context.Context, engine entity.ISSMoexEngine,
	market entity.ISSMoexMarket, board entity.ISSMoexBoard, ticker string) (ISSFullSecurityInfo, error) {
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

	result.LotSize = lotSize
	result.PriceDecimals = utils.SignsAfterDot(prevPrice)

	return result, nil
}

func (api *IssClient) GetStocksCurrentPrices(ctx context.Context, market entity.ISSMoexMarket,
	tickerInfos map[string]entity.ISSMoexBoard) (map[string]float64, error) {
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

	fmt.Println("data", data)

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
			return nil, fmt.Errorf("%s: failed to cast price from issreponse", op)
		}
		v, ok := tickerInfos[ticker]
		if ok && string(v) == board {
			m[ticker] = price
		}
	}

	return m, nil
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
