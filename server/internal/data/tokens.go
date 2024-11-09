package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/google/uuid"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string
	Model     repository.Token
}

func (t Token) New(userId uuid.UUID, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		Model: repository.Token{
			UserID: userId,
			Expiry: time.Now().Add(ttl),
			Scope:  scope,
		},
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Model.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, pt string) {
	v.CheckStrNotEmpty(pt, "token")
	v.Check(len(pt) == 26, "token", "must be 26 bytes long")
}
