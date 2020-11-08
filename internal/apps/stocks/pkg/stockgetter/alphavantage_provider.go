package stockgetter

import (
	"context"
	"fmt"
	"sort"

	"stockplay/pkg/alphavantage"
)

const (
	TimeModeIntraday = iota
	TimeModeDaily
	TimeModeWeekly
	TimeModeMonthly

	TimeInterval1Min = iota + 50
	TimeInterval5Min
	TimeInterval15Min
	TimeInterval30Min
	TimeInterval60Min
)

type Point struct {
	CurrentValue float64 `json:"current_value"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	Variation    float64 `json:"variation"`
	PrevClose    float64 `json:"previous_close"`
	Open         float64 `json:"open"`
	Volume       int64   `json:"volume"`
	Time         int64   `json:"time"`
}

type Stock struct {
	Points    []Point `json:"points"`
	MarketCap float64 `json:"market_cap"`
	AvgVolume int64   `json:"avg_volume"`
}

type AlphaVantageClient interface {
	GetStockTimeSeries(ctx context.Context, args alphavantage.GetStockArgs) ([]alphavantage.Stock, error)
}

type AlphaVantageStockGetter struct {
	client AlphaVantageClient
}

func NewAlphaVantageStockGetter(client AlphaVantageClient) *AlphaVantageStockGetter {
	return &AlphaVantageStockGetter{client: client}
}

type GetStockArgs struct {
	Mode     int
	Interval int
	Symbol   string
}

func (a *AlphaVantageStockGetter) Get(ctx context.Context, args GetStockArgs) (Stock, error) {
	resp, err := a.client.GetStockTimeSeries(ctx, alphavantage.GetStockArgs{
		Mode:     toAVFunc(args.Mode),
		Interval: toAVInterval(args.Interval),
		Symbol:   args.Symbol,
	})
	if err != nil {
		return Stock{}, fmt.Errorf("failed to get stock: %w", err)
	}

	// before parsing we make sure to sort by date ascending
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Date.Before(resp[j].Date)
	})

	return parseAVStocks(resp), nil
}

func parseAVStocks(stocks []alphavantage.Stock) Stock {
	var totalVolume int64
	var marketCap float64
	var res Stock
	var prevClose float64
	for _, s := range stocks {
		res.Points = append(res.Points, Point{
			CurrentValue: s.Close,
			Bid:          s.High,
			Ask:          s.Low,
			Variation:    s.High - s.Low,
			PrevClose:    prevClose,
			Open:         s.Open,
			Volume:       s.Volume,
			Time:         s.Date.Unix(),
		})

		totalVolume += s.Volume
		marketCap += s.Close * float64(s.Volume)
		prevClose = s.Close
	}

	res.AvgVolume = totalVolume / int64(len(stocks))
	res.MarketCap = marketCap

	return res
}

func toAVInterval(mode int) string {
	switch mode {
	case TimeInterval1Min:
		return alphavantage.Interval1min
	case TimeInterval15Min:
		return alphavantage.Interval15min
	case TimeInterval30Min:
		return alphavantage.Interval30min
	case TimeInterval60Min:
		return alphavantage.Interval60min
	}

	return alphavantage.Interval5min
}

func toAVFunc(mode int) string {
	switch mode {
	case TimeModeDaily:
		return alphavantage.ModeTimeSeriesDaily
	case TimeModeIntraday:
		return alphavantage.ModeTimeSeriesIntraday
	case TimeModeMonthly:
		return alphavantage.ModeTimeSeriesMonthly
	}

	return alphavantage.ModeTimeSeriesWeekly
}
