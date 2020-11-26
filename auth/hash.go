package auth

import (
	"golang.org/x/crypto/scrypt"
)

const PW_HASH_BYTES = 64

func CreateHashWithSalt(password string, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
