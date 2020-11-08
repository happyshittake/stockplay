package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	var tts = []struct {
		caseName     string
		handler      http.HandlerFunc
		expectedResp []byte
		expectedErr  error
	}{
		{
			caseName: "when error response from server",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("error"))
			},
			expectedResp: nil,
			expectedErr:  ErrServerError,
		},
		{
			caseName: "when success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("abcd1234"))
			},
			expectedResp: []byte("abcd1234"),
			expectedErr:  nil,
		},
	}

	for idx, tt := range tts {
		logTestcase := fmt.Sprintf("[TESTCASE %d]", idx)
		t.Log(logTestcase, tt.caseName)

		server := httptest.NewServer(tt.handler)

		c := Client{
			httpClient: http.DefaultClient,
			host:       server.URL,
		}

		resp, err := c.Encrypt(context.Background(), []byte("1234"))
		if err != nil {
			if !errors.Is(err, tt.expectedErr) {
				t.Error("expected err:", tt.expectedErr, ", is not err:", err)
			}
		}

		if string(resp) != string(tt.expectedResp) {
			t.Error("expected resp:", string(tt.expectedResp), ", not equal:", string(resp))
		}

		server.Close()
	}
}
