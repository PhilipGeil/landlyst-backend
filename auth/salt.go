package auth

import (
	"crypto/rand"
	"io"
)

const PW_SALT_BYTES = 32

func CreateSalt() ([]byte, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
