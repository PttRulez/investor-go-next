package iss_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ISSecurityInfo struct {
	Board     entity.ISSMoexBoard
	Engine    entity.ISSMoexEngine
	Market    entity.ISSMoexMarket
	Name      string
	ShortName string
}

func (api *IssClient) GetSecurityInfoBySecid(secid string) (*ISSecurityInfo, error) {
	uri := fmt.Sprintf("%s/securities/%s.json", api.baseUrl, secid)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssClient.GetSecurityByTicker http_controllers.NewRequest]: %w", err)
	}

	// фильтруем только то что нам нужно
	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine,is_primary")
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

	var (
		name, shortname string
		board           entity.ISSMoexBoard
		market          entity.ISSMoexMarket
		engine          entity.ISSMoexEngine
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
		// если is_primary
		if item[3].(float64) == 1 {
			boardData = item
		}
	}
	board = entity.ISSMoexBoard(boardData[0].(string))
	market = entity.ISSMoexMarket(boardData[1].(string))
	engine = entity.ISSMoexEngine(boardData[2].(string))

	return &ISSecurityInfo{
		Board:     board,
		Engine:    engine,
		Market:    market,
		Name:      name,
		ShortName: shortname,
	}, nil
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
