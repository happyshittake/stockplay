package alphavantage

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

const (
	ModeTimeSeriesIntraday = "TIME_SERIES_INTRADAY"
	ModeTimeSeriesDaily    = "TIME_SERIES_DAILY"
	ModeTimeSeriesWeekly   = "TIME_SERIES_WEEKLY"
	ModeTimeSeriesMonthly  = "TIME_SERIES_MONTHLY"

	Interval1min  = "1min"
	Interval5min  = "5min"
	Interval15min = "15min"
	Interval30min = "30min"
	Interval60min = "60min"

	layoutIntraday = "2006-01-02 15:04:05"
	layoutStd      = "2006-01-02"
)

var (
	ErrServerResponse = errors.New("server response error")
)

type Client struct {
	httpClient *http.Client
	host       string
	apiKey     string
}

func NewClient(httpClient *http.Client, host, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		host:       host,
		apiKey:     apiKey,
	}
}

type Stock struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
	Date   time.Time
}

type GetStockArgs struct {
	Mode     string
	Interval string
	Symbol   string
}

func (c *Client) GetStockTimeSeries(ctx context.Context, args GetStockArgs) ([]Stock, error) {
	q := url.Values{}
	q.Set("function", args.Mode)
	q.Set("symbol", args.Symbol)
	q.Set("apikey", c.apiKey)
	q.Set("datatype", "csv")
	if args.Mode == ModeTimeSeriesIntraday {
		q.Set("interval", args.Interval)
	}

	urlpath := c.host + "/query?" + q.Encode()

	req, err := http.NewRequest(http.MethodGet, urlpath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dumpedResponse, _ := httputil.DumpResponse(resp, true)

		return nil, fmt.Errorf("got error response \n%s\n: %w", string(dumpedResponse), ErrServerResponse)
	}

	return parseBody(args.Mode, resp.Body)
}

func parseBody(mode string, body io.Reader) ([]Stock, error) {
	csvReader := csv.NewReader(body)
	csvReader.LazyQuotes = true

	csvReader.Read() // skip header

	layout := layoutStd
	if mode == ModeTimeSeriesIntraday {
		layout = layoutIntraday
	}

	var stocks []Stock
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return stocks, nil
			}

			return stocks, fmt.Errorf("error reading csv row: %w", err)
		}

		date, err := time.Parse(layout, row[0])
		if err != nil {
			return stocks, fmt.Errorf("failed to parse timestamp %s: %w", row[0], err)
		}

		open, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return stocks, fmt.Errorf("failed to parse open field %s: %w", row[1], err)
		}

		high, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return stocks, fmt.Errorf("failed to parse high field %s: %w", row[2], err)
		}

		low, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return stocks, fmt.Errorf("failed to parse low field %s: %w", row[3], err)
		}

		cls, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return stocks, fmt.Errorf("failed to parse close field %s: %w", row[4], err)
		}

		volume, err := strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			return stocks, fmt.Errorf("failed to parse open field %s: %w", row[5], err)
		}

		stocks = append(stocks, Stock{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  cls,
			Volume: volume,
			Date:   date,
		})
	}
}
