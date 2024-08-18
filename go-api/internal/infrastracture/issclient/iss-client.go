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
)

type ISSecurityInfo struct {
	Board     entity.ISSMoexBoard
	Engine    entity.ISSMoexEngine
	Market    entity.ISSMoexMarket
	Name      string
	ShortName string
	Secid     string
	// Только для облигаций
	CouponPercent   float64
	CouponValue     float64
	CouponFrequency int       // частота выплаты купонов в год
	IssueDate       time.Time // облигации
	FaceValue       int       // номинальная стоимость
	MatDate         time.Time // дата погашения облиги

}

func (api *IssClient) GetSecurityInfoBySecid(ctx context.Context, secid string) (*ISSecurityInfo, error) {
	const op = "issclient.GetSecurityInfoBySecid"

	uri := fmt.Sprintf("%s/securities/%s.json", api.baseURL, secid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	params := url.Values{}
	// фильтруем только то что нам нужно
	params.Add("iss.meta", "off")
	params.Add("description.columns", "name,value")
	params.Add("boards.columns", "boardid,market,engine")
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

	var data MoexAPIResponseSecurityInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
			freq, err := strconv.Atoi(item[1])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			result.CouponFrequency = freq
		case "COUPONPERCENT":
			percent, err := strconv.ParseFloat(item[1], 32)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			result.CouponPercent = percent
		case "COUPONVALUE":
			percent, err := strconv.ParseFloat(item[1], 32)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			result.CouponPercent = percent
		case "ISSUEDATE":
			t, err := time.Parse("2006-01-02", item[1])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			result.IssueDate = t
		case "MATDATE":
			t, err := time.Parse("2006-01-02", item[1])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			result.MatDate = t
		case "FACEVALUE":
			faceValue, err := strconv.Atoi(item[1])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
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
	const op = "issclient.GetStocksCurrentPrices"

	uri := fmt.Sprintf("%s/engines/stock/markets/%s/securities.json", api.baseURL, market)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	params := url.Values{}
	params.Add("iss.meta", "off")
	params.Add("securities", url.QueryEscape(strings.Join(tickers, ",")))
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

	var m Prices
	for _, i := range data.Securities.Data {
		ticker, ok := i[0].(string)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast ticker from issreponse", op)
		}

		board, ok := i[1].(entity.ISSMoexBoard)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast board from issreponse", op)
		}

		price, ok := i[2].(float64)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast price from issreponse", op)
		}

		m[ticker][board] = price
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
