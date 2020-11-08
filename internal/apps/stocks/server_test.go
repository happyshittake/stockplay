package stocks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"stockplay/internal/apps/stocks/pkg/stockgetter"
)

type failStockGetter int

func (f failStockGetter) Get(ctx context.Context, args stockgetter.GetStockArgs) (stockgetter.Stock, error) {
	return stockgetter.Stock{}, errors.New("any err")
}

type successStockGetter int

func (s successStockGetter) Get(ctx context.Context, args stockgetter.GetStockArgs) (stockgetter.Stock, error) {
	return stockgetter.Stock{}, nil
}

type failEncSvc int

func (f failEncSvc) Encrypt(ctx context.Context, text []byte) ([]byte, error) {
	return nil, errors.New("any err")
}

type successEncSvc int

func (s successEncSvc) Encrypt(ctx context.Context, text []byte) ([]byte, error) {
	return []byte("abcd123"), nil
}

func TestServer_HandleGetStock(t *testing.T) {
	var tts = []struct {
		caseName           string
		enc                EncryptService
		sg                 StockGetter
		expectedStatusCode int
		expectedBody       string
		symbol             string
	}{
		{
			caseName:           "when error getting the stock",
			enc:                failEncSvc(1),
			sg:                 failStockGetter(1),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "internal error",
			symbol:             "abcd123",
		},
		{
			caseName:           "when error encrypting data",
			enc:                failEncSvc(1),
			sg:                 successStockGetter(1),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "internal error",
			symbol:             "abcd123",
		},
		{
			caseName:           "when success",
			enc:                successEncSvc(1),
			sg:                 successStockGetter(1),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "abcd123",
			symbol:             "abcd123",
		},
	}

	for idx, tt := range tts {
		logTestcase := fmt.Sprintf("[TESTCASE %d]", idx)
		t.Log(logTestcase, tt.caseName)

		s := Server{
			stockGetter: tt.sg,
			encService:  tt.enc,
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(logTestcase, err)
		}

		rw := httptest.NewRecorder()

		s.HandleGetStock().ServeHTTP(rw, req)

		if rw.Code != tt.expectedStatusCode {
			t.Errorf("%s status code [%d] not equal expected [%d]", logTestcase, rw.Code, tt.expectedStatusCode)
		}

		if rw.Body.String() != tt.expectedBody {
			t.Errorf("%s body [%s] not equal expected [%s]", logTestcase, rw.Body.String(), tt.expectedBody)
		}
	}
}
