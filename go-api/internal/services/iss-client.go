package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

type IssApiService struct {
	baseUrl string
	client  *http.Client
}

func NewIssApiService() *IssApiService {
	return &IssApiService{
		baseUrl: "https://iss.moex.com/iss",
		client:  &http.Client{},
	}
}

func (api *IssApiService) GetSecurityInfoByTicker(ticker string) (*tmoex.Security, error) {
	uri := fmt.Sprintf("%s/securities/%s.json", api.baseUrl, ticker)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker http.NewRequest]: %w", err)
	}

	// фильтруем только то что нам нужно
	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine,is_primary")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker api.client.Do(req)]: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker ReadAll(body)]: %w", err)
	}

	data := &types.MoexApiResponseSecurityInfo{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker json.Unmarshal(body)]: %w", err)
	}

	var (
		name, shortname string
		board           tmoex.Board
		market          tmoex.Market
		engine          tmoex.Engine
		// ok              bool
	)

	for _, item := range data.Description.Data {
		if item[0] == "NAME" {
			name = item[1]
		}
		if item[0] == "SHORTNAME" {
			shortname = item[1]
		}
	}

	var boardData [4]any
	for _, item := range data.Boards.Data {
		if item[3].(float64) == 1 {
			boardData = item
		}
	}
	board = tmoex.Board(boardData[0].(string))
	market = tmoex.Market(boardData[1].(string))
	engine = tmoex.Engine(boardData[2].(string))

	return &tmoex.Security{
		Board:     board,
		Engine:    engine,
		Market:    market,
		Name:      name,
		ShortName: shortname,
		Ticker:    ticker,
	}, nil
}

func (api *IssApiService) GetStocksCurrentPrices(ctx context.Context, market tmoex.Market, tickers []string) (*types.MoexApiResponseCurrentPrices, error) {
	uri := fmt.Sprintf("%s/engines/stock/markets/%s/securities.json", api.baseUrl, market)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices http.NewRequest]: %w", err)
	}

	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("securities", url.QueryEscape(strings.Join(tickers, ",")))
	params.Add("securities.columns", "SECID,BOARDID,PREVPRICE")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices client.Do(req)]: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices ReadAll(body)]: %w", err)
	}

	data := &types.MoexApiResponseCurrentPrices{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices json.Unmarshal(body)]: %w", err)
	}

	return data, nil
}
