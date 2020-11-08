package encryptor

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type successEnc int

func (s successEnc) Encrypt(text string) (string, error) {
	return "abcd1234", nil
}

type failEnc int

func (f failEnc) Encrypt(text string) (string, error) {
	return "", errors.New("any error")
}

func TestEncryptionHandler(t *testing.T) {
	var tts = []struct {
		caseName           string
		enc                Encryptor
		expectedStatusCode int
		expectedBody       string
		requestBody        []byte
	}{
		{
			caseName:           "when failed to encrypt request body",
			enc:                failEnc(1),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "failed to encrypt message",
			requestBody:        []byte("any body"),
		},
		{
			caseName:           "success case",
			enc:                successEnc(1),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "abcd1234",
			requestBody:        []byte("any body"),
		},
	}

	for idx, tt := range tts {
		logTestcase := fmt.Sprintf("[TESTCASE %d]", idx)
		t.Log(logTestcase, tt.caseName)

		s := Server{enc: tt.enc}

		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.requestBody))
		if err != nil {
			t.Error(logTestcase, err)
		}

		rw := httptest.NewRecorder()

		s.HandleEncrypt().ServeHTTP(rw, req)

		if rw.Result().StatusCode != tt.expectedStatusCode {
			t.Errorf("%s expected status code [%d] not equal to received status code [%d]", logTestcase, tt.expectedStatusCode, rw.Result().StatusCode)
		}

		if rw.Body.String() != tt.expectedBody {
			t.Errorf("%s expected body [%s] not equal to received body [%s]", logTestcase, tt.expectedBody, rw.Body.String())
		}
	}
}
