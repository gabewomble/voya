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
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash

	return nil
}

func (p password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
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

func (u UserInput) validateEmail(v *validator.Validator) {
	v.CheckStrNotEmpty(u.Email, "email")
	v.Check(validator.Matches(u.Email, validator.EmailRX), "email", "must be a valid email address")
}

func (u UserInput) validatePasswordPlaintext(v *validator.Validator) {
	v.CheckStrNotEmpty(*u.Password.plaintext, "password")
	v.Check(len(*u.Password.plaintext) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(*u.Password.plaintext) <= 72, "password", "must not be more than 72 bytes long")
}

func (u UserInput) validateName(v *validator.Validator) {
	v.CheckStrNotEmpty(u.Name, "name")
	v.Check(len(u.Name) <= 500, "name", "must not be more than 500 bytes long")
}

func (u UserInput) Validate(v *validator.Validator) {
	u.validateName(v)
	u.validateEmail(v)

	if u.Password.plaintext != nil {
		u.validatePasswordPlaintext(v)
		if !v.Valid() {
			return
		}
	}

	if u.Password.hash == nil {
		panic("missing password hash for user")
	}
}
