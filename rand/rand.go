package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes : didn't read enough random bytes")
	}

	return b, nil
}

func String(numberOfBytes int) (string, error) {
	b, err := Bytes(numberOfBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
