package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"server/internal/repository"
	"server/internal/validator"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tokenScope struct {
	Activation     string
	Authentication string
}

var TokenScope = tokenScope{
	Activation:     "activation",
	Authentication: "authentication",
}

type tokenErrors struct {
	NotFound     error
	InvalidToken error
}

var TokenErr = tokenErrors{
	NotFound:     errors.New("token not found"),
	InvalidToken: errors.New("invalid token"),
}

type Token struct {
	Plaintext    string
	RefreshToken string
	Model        repository.Token
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

	// Generate refresh token
	refreshBytes := make([]byte, 16)
	_, err = rand.Read(refreshBytes)
	if err != nil {
		return nil, err
	}

	token.RefreshToken = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(refreshBytes)
	refreshHash := sha256.Sum256([]byte(token.RefreshToken))
	token.Model.RefreshToken = refreshHash[:]

	return token, nil
}

func GetTokenPlainTextFromContext(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", TokenErr.NotFound
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", TokenErr.InvalidToken
	}

	token := headerParts[1]

	v := validator.New()

	if ValidateTokenPlaintext(v, token, "token"); !v.Valid() {
		return "", TokenErr.InvalidToken
	}

	return token, nil
}

func GetTokenHash(plaintext string) [32]byte {
	return sha256.Sum256([]byte(plaintext))
}

func ValidateTokenPlaintext(v *validator.Validator, pt string, fieldName string) {
	v.CheckStrNotEmpty(pt, fieldName)
	v.Check(len(pt) == 26, fieldName, "must be 26 bytes long")
}
