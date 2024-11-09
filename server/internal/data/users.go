package data

import (
	"errors"
	"server/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name     string
	Email    string
	Password password
}

type password struct {
	plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.Hash = hash

	return nil
}

func CompareUserHashAndPassword(hash []byte, pt string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pt))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateUserEmail(v *validator.Validator, email string) {
	v.CheckStrNotEmpty(email, "email")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUserPasswordPlaintext(v *validator.Validator, pt string) {
	v.CheckStrNotEmpty(pt, "password")
	v.Check(len(pt) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(pt) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUserName(v *validator.Validator, name string) {
	v.CheckStrNotEmpty(name, "name")
	v.Check(len(name) <= 500, "name", "must not be more than 500 bytes long")
}

func (u UserInput) Validate(v *validator.Validator) {
	ValidateUserName(v, u.Name)
	ValidateUserEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidateTokenPlaintext(v, *u.Password.plaintext)
		if !v.Valid() {
			return
		}
	}

	if u.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
