package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"stockplay/internal/apps/encryptor/pkg/client"
	"stockplay/internal/apps/stocks"
	"stockplay/internal/apps/stocks/pkg/stockgetter"
	"stockplay/pkg/alphavantage"
)

func main() {

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	encClient := client.NewClient(httpClient, os.Getenv("ENCRYPTOR_HOST"))
	alphaVantageClient := alphavantage.NewClient(
		httpClient,
		os.Getenv("ALPHAVANTAGE_HOST"),
		os.Getenv("ALPHAVANTAGE_KEY"),
	)
	stockGetter := stockgetter.NewAlphaVantageStockGetter(alphaVantageClient)

	handler := stocks.NewServer(stockGetter, encClient)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      handler.HandleGetStock(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("starting stock service at", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
