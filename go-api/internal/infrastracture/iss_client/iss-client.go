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

func (api *IssApiService) GetSecurityInfoBySecid(secid string) (*ISSecurityInfo, error) {
	uri := fmt.Sprintf("%s/securities/%s.json", api.baseUrl, secid)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker http_controllers.NewRequest]: %w", err)
	}

	// фильтруем только то что нам нужно
	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine,is_primary")
	req.URL.RawQuery = params.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker controller.client.Do(req)]: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker ReadAll(body)]: %w", err)
	}

	data := &MoexApiResponseSecurityInfo{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetSecurityByTicker json.Unmarshal(body)]: %w", err)
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

func (api *IssApiService) GetStocksCurrentPrices(ctx context.Context, market entity.ISSMoexMarket,
	tickers []string) (*MoexApiResponseCurrentPrices, error) {
	uri := fmt.Sprintf("%s/engines/stock/markets/%s/securities.json", api.baseUrl, market)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices http_controllers.NewRequest]: %w", err)
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

	data := &MoexApiResponseCurrentPrices{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("[IssApiService.GetStocksCurrentPrices json.Unmarshal(body)]: %w", err)
	}

	return data, nil
}

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
