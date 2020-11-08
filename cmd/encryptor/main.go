package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"stockplay/internal/apps/encryptor"
)

func main() {

	key := []byte(os.Getenv("ENCRYPTOR_KEY"))
	enc, err := encryptor.NewAes256Encryption(key)
	if err != nil {
		log.Fatal("failed to create encryption method ", len(key), err)
	}

	handler := encryptor.NewServer(enc)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      handler.HandleEncrypt(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("starting encryptor service at ", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
