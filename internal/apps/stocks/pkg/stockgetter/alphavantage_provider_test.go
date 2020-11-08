package stockgetter

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"stockplay/pkg/alphavantage"
)

var now = time.Now()

type failAvClient int

func (f failAvClient) GetStockTimeSeries(ctx context.Context, args alphavantage.GetStockArgs) ([]alphavantage.Stock, error) {
	return nil, errors.New("any error")
}

type successAvClient int

func (s successAvClient) GetStockTimeSeries(ctx context.Context, args alphavantage.GetStockArgs) ([]alphavantage.Stock, error) {
	return []alphavantage.Stock{
		{
			Open:   100.00,
			High:   90.00,
			Low:    10.00,
			Close:  80.00,
			Volume: 100,
			Date:   now,
		},
		{
			Open:   200.00,
			High:   100.00,
			Low:    20.00,
			Close:  90.00,
			Volume: 110,
			Date:   now.Add(1 * time.Hour),
		},
	}, nil
}

func TestAlphaVantageStockGetter_Get(t *testing.T) {
	var tts = []struct {
		caseName     string
		mode         int
		symbol       string
		interval     int
		client       AlphaVantageClient
		expectedResp Stock
		expectedErr  bool
	}{
		{
			caseName:     "when response from client error",
			mode:         TimeModeMonthly,
			symbol:       "abcd123",
			interval:     TimeInterval30Min,
			client:       failAvClient(1),
			expectedResp: Stock{},
			expectedErr:  true,
		},
		{
			caseName: "when success",
			mode:     TimeModeMonthly,
			symbol:   "abcd123",
			interval: TimeInterval30Min,
			client:   successAvClient(1),
			expectedResp: Stock{
				Points: []Point{
					{
						CurrentValue: 80.00,
						Bid:          90.00,
						Ask:          10.00,
						Variation:    90.00 - 10.00,
						PrevClose:    0,
						Open:         100.00,
						Volume:       100,
						Time:         now.Unix(),
					},
					{
						CurrentValue: 90.00,
						Bid:          100.00,
						Ask:          20.00,
						Variation:    100.00 - 20.00,
						PrevClose:    80.00,
						Open:         200.00,
						Volume:       110,
						Time:         now.Add(1 * time.Hour).Unix(),
					},
				},
				MarketCap: (80.00 * 100.00) + (90.00 * 110.00),
				AvgVolume: (100 + 110) / 2,
			},
			expectedErr: false,
		},
	}

	for idx, tt := range tts {
		logTestcase := fmt.Sprintf("[TESTCASE %d]", idx)
		t.Log(logTestcase, tt.caseName)

		c := AlphaVantageStockGetter{
			client: tt.client,
		}

		resp, err := c.Get(context.Background(), GetStockArgs{
			Mode:     tt.mode,
			Interval: tt.interval,
			Symbol:   tt.symbol,
		})
		if err != nil {
			if !tt.expectedErr {
				t.Error(logTestcase, "unexpected err", err)
			}
		}

		if resp.MarketCap != tt.expectedResp.MarketCap && resp.AvgVolume != tt.expectedResp.AvgVolume {
			t.Errorf("%s received value %+v not equal expected value %+v", logTestcase, resp, tt.expectedResp)
		}

		if len(resp.Points) != len(tt.expectedResp.Points) {
			t.Errorf("%s received value %+v not equal expected value %+v", logTestcase, resp, tt.expectedResp)
		}

		for k, point := range resp.Points {
			if !reflect.DeepEqual(point, tt.expectedResp.Points[k]) {
				t.Errorf("%s received value %+v not equal expected value %+v", logTestcase, point, tt.expectedResp.Points[k])
			}
		}
	}
}
