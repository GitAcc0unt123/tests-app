package hash

import (
	"crypto/sha512"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type SHA512Hasher struct {
	salt string
}

func NewSHA512Hasher(salt string) *SHA512Hasher {
	return &SHA512Hasher{salt: salt}
}

func (h *SHA512Hasher) Hash(password string) (string, error) {
	hash := sha512.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
