package encryptor

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Encryptor interface {
	Encrypt(text string) (string, error)
}

type Server struct {
	enc Encryptor
}

func NewServer(enc Encryptor) *Server {
	return &Server{enc: enc}
}

func (s *Server) HandleEncrypt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read body", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to read body"))
			return
		}

		encrypted, err := s.enc.Encrypt(string(body))
		if err != nil {
			log.Println("failed to encrypt message", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to encrypt message"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encrypted))
	}
}
