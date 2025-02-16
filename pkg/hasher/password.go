package hasher

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type BcryptHasher struct {
	salt string
}

func New(salt string) *BcryptHasher {
	return &BcryptHasher{
		salt: salt,
	}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	salt := []byte(h.salt)
	rand.Read(salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	hashedPassword, err := bcrypt.GenerateFromPassword(append(salt, []byte(password)...), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("$2b$10$%s%s", saltBase64, base64.StdEncoding.EncodeToString(hashedPassword))

	return hash, nil
}

func (h *BcryptHasher) Verify(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 {
		return false, fmt.Errorf("invalid hash format")
	}

	saltBase64 := parts[2]
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), append(salt, []byte(password)...))
	return err == nil, nil
}
