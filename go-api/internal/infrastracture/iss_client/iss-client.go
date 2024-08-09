package iss_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ISSecurityInfo struct {
	Board     entity.ISSMoexBoard
	Engine    entity.ISSMoexEngine
	Market    entity.ISSMoexMarket
	Name      string
	ShortName string
	Secid     string
	// Только для облигаций
	CouponPercent   float32
	CouponValue     float32
	CouponFrequency int8      // частота выплаты купонов в год
	IssueDate       time.Time // облигации
	FaceValue       int       // номинальная стоимость
	MatDate         time.Time // дата погашения облиги

}

func (api *IssClient) GetSecurityInfoBySecid(secid string) (*ISSecurityInfo, error) {
	uri := fmt.Sprintf("%s/securities/%s.json", api.baseUrl, secid)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetSecurityByTicker http_controllers.NewRequest]: %w", err)
	}

	params := url.Values{}
	// фильтруем только то что нам нужно
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetSecurityByTicker controller.client.Do(req)]: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetSecurityByTicker ReadAll(body)]: %w", err)
	}

	data := &MoexApiResponseSecurityInfo{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetSecurityByTicker json.Unmarshal(body)]: %w", err)
	}

	result := ISSecurityInfo{}
	for _, item := range data.Description.Data {
		switch item[0] {
		case "NAME":
			result.Name = item[1]
		case "SHORTNAME":
			result.ShortName = item[1]
		// Только для облигаций:
		case "COUPONFREQUENCY":
			freq, err := strconv.ParseInt(item[1], 10, 8)
			if err != nil {
				return nil, err
			}
			result.CouponFrequency = int8(freq)
		case "COUPONPERCENT":
			percent, err := strconv.ParseFloat(item[1], 32)
			if err != nil {
				return nil, err
			}
			result.CouponPercent = float32(percent)
		case "COUPONVALUE":
			percent, err := strconv.ParseFloat(item[1], 32)
			if err != nil {
				return nil, err
			}
			result.CouponPercent = float32(percent)
		case "ISSUEDATE":
			t, err := time.Parse("2006-01-02", item[1])
			if err != nil {
				return nil, err
			}
			result.IssueDate = t
		case "MATDATE":
			t, err := time.Parse("2006-01-02", item[1])
			if err != nil {
				return nil, err
			}
			result.MatDate = t
		case "FACEVALUE":
			faceValue, err := strconv.Atoi(item[1])
			if err != nil {
				return nil, err
			}
			result.FaceValue = faceValue
		}
	}

	boardData := data.Boards.Data[0]
	result.Board = entity.ISSMoexBoard(boardData[0])
	result.Market = entity.ISSMoexMarket(boardData[1])
	result.Engine = entity.ISSMoexEngine(boardData[2])

	return &result, nil
}

type Prices map[string]map[entity.ISSMoexBoard]float64

func (api *IssClient) GetStocksCurrentPrices(ctx context.Context, market entity.ISSMoexMarket,
	tickers []string) (Prices, error) {
	uri := fmt.Sprintf("%s/engines/stock/markets/%s/securities.json", api.baseUrl, market)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices http_controllers.NewRequest]: %w", err)
	}

	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("securities", url.QueryEscape(strings.Join(tickers, ",")))
	params.Add("securities.columns", "SECID,BOARDID,PREVPRICE")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices client.Do(req)]: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices ReadAll(body)]: %w", err)
	}

	data := &MoexApiResponseCurrentPrices{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices json.Unmarshal(body)]: %w", err)
	}

	var m Prices
	for _, i := range data.Securities.Data {
		ticker, ok := i[0].(string)
		if !ok {
			return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices ticker]: %w", err)
		}

		board, ok := i[1].(entity.ISSMoexBoard)
		if !ok {
			return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices ticker]: %w", err)
		}

		price, ok := i[2].(float64)
		if !ok {
			return nil, fmt.Errorf("[IssClient.GetStocksCurrentPrices ticker]: %w", err)
		}

		m[ticker][board] = price
	}
	return m, nil
}

type IssClient struct {
	baseUrl string
	client  *http.Client
}

func NewISSClient() *IssClient {
	return &IssClient{
		baseUrl: "https://iss.moex.com/iss",
		client:  &http.Client{},
	}
}
