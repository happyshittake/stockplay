package encryptor

import (
	"testing"
)

func TestEncDec(t *testing.T) {
	var tts = []struct {
		caseName  string
		key       string
		plainText string
		err       error
	}{
		{
			caseName:  "key length is not 32",
			key:       "abc123",
			plainText: "this is a plain text",
			err:       ErrInvalidKey,
		},
		{
			caseName:  "successfully encrypt then decrypt it",
			key:       "abcdefghijklmnopqrstuvwxyz012345",
			plainText: "valid plain text",
			err:       nil,
		},
	}

	for _, tt := range tts {
		t.Log(tt.caseName)

		encryptor, err := NewAes256Encryption([]byte(tt.key))
		if !checkErr(t, tt.err, err) {
			continue
		}

		encrypted, err := encryptor.Encrypt(tt.plainText)
		if !checkErr(t, tt.err, err) {
			continue
		}

		decrypted, err := encryptor.Decrypt(encrypted)
		if !checkErr(t, tt.err, err) {
			continue
		}

		if decrypted != tt.plainText {
			t.Error(
				"plaintext is not the same compared to decrypted text",
				"decrypted", decrypted,
				"plain text", tt.plainText,
			)
		}

	}
}

func checkErr(t *testing.T, expectedErr error, err error) bool {
	if err != nil {
		if expectedErr == err {
			return false
		}

		t.Error("test fail", err)
		return false
	}

	return true
}
