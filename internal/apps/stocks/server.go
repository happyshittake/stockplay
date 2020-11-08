package stocks

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"stockplay/internal/apps/stocks/pkg/stockgetter"
)

type StockGetter interface {
	Get(ctx context.Context, args stockgetter.GetStockArgs) (stockgetter.Stock, error)
}

type EncryptService interface {
	Encrypt(ctx context.Context, text []byte) ([]byte, error)
}

type Server struct {
	stockGetter StockGetter
	encService  EncryptService
}

func NewServer(stockGetter StockGetter, encService EncryptService) *Server {
	return &Server{
		stockGetter: stockGetter,
		encService:  encService,
	}
}

func (s *Server) HandleGetStock() http.HandlerFunc {
	const (
		defaultMode     = stockgetter.TimeModeWeekly
		defaultInterval = stockgetter.TimeInterval60Min
	)

	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		mode := defaultMode
		if q.Get("mode") != "" {
			mode, _ = strconv.Atoi(q.Get("mode"))
		}

		interval := defaultInterval
		if q.Get("interval") != "" {
			interval, _ = strconv.Atoi(q.Get("interval"))
		}

		resp, err := s.stockGetter.Get(r.Context(), stockgetter.GetStockArgs{
			Mode:     mode,
			Interval: interval,
			Symbol:   q.Get("symbol"),
		})
		if err != nil {
			log.Println("got error when getting stock data", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}

		pretty, _ := json.MarshalIndent(resp, "", "  ")
		log.Printf("stock data \n%s", string(pretty))

		text, err := json.Marshal(resp)
		if err != nil {
			log.Println("got error when marshalling data", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}

		encrypted, err := s.encService.Encrypt(r.Context(), text)
		if err != nil {
			log.Println("got error when encrypting data", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(encrypted)
	}
}
