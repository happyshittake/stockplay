package alphavantage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetStockTimeSeries(t *testing.T) {
	var tts = []struct {
		caseName     string
		mode         string
		symbol       string
		interval     string
		handler      func(logtag string, t *testing.T) http.HandlerFunc
		expectedResp []Stock
		expectedErr  error
	}{
		{
			caseName: "when error response from server",
			mode:     ModeTimeSeriesIntraday,
			symbol:   "abcde",
			interval: Interval5min,
			handler: func(logtag string, t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					q := r.URL.Query()

					if q.Get("function") != ModeTimeSeriesIntraday {
						t.Errorf("%s mode [%s] is not equal [%s]", logtag, q.Get("function"), ModeTimeSeriesIntraday)
					}

					if q.Get("symbol") != "abcde" {
						t.Errorf("%s symbol [%s] is not equal [%s]", logtag, q.Get("symbol"), "abcde")
					}

					if q.Get("interval") != Interval5min {
						t.Errorf("%s interval [%s] is not equal [%s]", logtag, q.Get("interval"), Interval5min)
					}

					if q.Get("apikey") != "demo" {
						t.Errorf("%s demo [%s] is not equal [%s]", logtag, q.Get("apikey"), "demo")
					}

					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("error"))
				}
			},
			expectedResp: nil,
			expectedErr:  ErrServerResponse,
		},
		{
			caseName: "when success intraday",
			mode:     ModeTimeSeriesIntraday,
			interval: Interval5min,
			symbol:   "abcde",
			handler: func(logtag string, t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					q := r.URL.Query()

					if q.Get("function") != ModeTimeSeriesIntraday {
						t.Errorf("%s mode [%s] is not equal [%s]", logtag, q.Get("function"), ModeTimeSeriesIntraday)
					}

					if q.Get("symbol") != "abcde" {
						t.Errorf("%s symbol [%s] is not equal [%s]", logtag, q.Get("symbol"), "abcde")
					}

					if q.Get("interval") != Interval5min {
						t.Errorf("%s interval [%s] is not equal [%s]", logtag, q.Get("interval"), Interval5min)
					}

					if q.Get("apikey") != "demo" {
						t.Errorf("%s demo [%s] is not equal [%s]", logtag, q.Get("apikey"), "demo")
					}

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`timestamp,open,high,low,close,volume
2020-11-06 20:00:00,114.4400,114.4400,114.4400,114.4400,457`))
				}
			},
			expectedResp: []Stock{
				{
					Open:   114.4400,
					High:   114.4400,
					Low:    114.4400,
					Close:  114.4400,
					Volume: 457,
					Date:   time.Date(2020, 11, 6, 20, 0, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
		{
			caseName: "when success other than intraday",
			mode:     ModeTimeSeriesMonthly,
			symbol:   "abcde",
			handler: func(logtag string, t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					q := r.URL.Query()

					if q.Get("function") != ModeTimeSeriesMonthly {
						t.Errorf("%s mode [%s] is not equal [%s]", logtag, q.Get("function"), ModeTimeSeriesMonthly)
					}

					if q.Get("symbol") != "abcde" {
						t.Errorf("%s symbol [%s] is not equal [%s]", logtag, q.Get("symbol"), "abcde")
					}

					if q.Get("apikey") != "demo" {
						t.Errorf("%s demo [%s] is not equal [%s]", logtag, q.Get("apikey"), "demo")
					}

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`timestamp,open,high,low,close,volume
2020-11-06,114.4400,114.4400,114.4400,114.4400,457`))
				}
			},
			expectedResp: []Stock{
				{
					Open:   114.4400,
					High:   114.4400,
					Low:    114.4400,
					Close:  114.4400,
					Volume: 457,
					Date:   time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
	}

	for idx, tt := range tts {
		logTestcase := fmt.Sprintf("[TESTCASE %d]", idx)
		t.Log(logTestcase, tt.caseName)

		srv := httptest.NewServer(tt.handler(logTestcase, t))

		c := Client{
			httpClient: http.DefaultClient,
			host:       srv.URL,
			apiKey:     "demo",
		}

		resp, err := c.GetStockTimeSeries(context.Background(), GetStockArgs{
			Mode:     tt.mode,
			Interval: tt.interval,
			Symbol:   tt.symbol,
		})
		if err != nil {
			if !errors.Is(err, tt.expectedErr) {
				t.Error(logTestcase, "expected err:", tt.expectedErr, ", is not err:", err)
			}
		}

		if len(resp) != len(tt.expectedResp) {
			t.Error(logTestcase, "expected len resp:", len(resp), ", not equal:", len(tt.expectedResp))
		}

		for k, resp := range resp {
			if !resp.Date.Equal(tt.expectedResp[k].Date) {
				t.Error(logTestcase, "expected date:", tt.expectedResp[k].Date, ", not equal:", resp.Date)
			}

			if resp.Open != tt.expectedResp[k].Open {
				t.Error(logTestcase, "expected open:", tt.expectedResp[k].Open, ", not equal:", resp.Open)
			}

			if resp.Close != tt.expectedResp[k].Close {
				t.Error(logTestcase, "expected close:", tt.expectedResp[k].Close, ", not equal:", resp.Close)
			}

			if resp.High != tt.expectedResp[k].High {
				t.Error(logTestcase, "expected high:", tt.expectedResp[k].High, ", not equal:", resp.High)
			}

			if resp.Low != tt.expectedResp[k].Low {
				t.Error(logTestcase, "expected low:", tt.expectedResp[k].Low, ", not equal:", resp.Low)
			}

			if resp.Volume != tt.expectedResp[k].Volume {
				t.Error(logTestcase, "expected volume:", tt.expectedResp[k].Volume, ", not equal:", resp.Volume)
			}
		}

		srv.Close()
	}
}
