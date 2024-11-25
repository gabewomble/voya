package data

import (
	"errors"
	"server/internal/repository"
	"server/internal/validator"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string
	Name     string
	Email    string
	Password password
}

type password struct {
	plaintext *string
	Hash      []byte
}

type userConstraints struct {
	UniqueEmail    string
	UniqueUsername string
}

var UserConstraint = userConstraints{
	UniqueEmail:    "users_email_key",
	UniqueUsername: "idx_users_username",
}

var AnonymousUser = &repository.User{}

func UserIsAnonymous(u *repository.User) bool {
	return u == AnonymousUser
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

func ValidateUsername(v *validator.Validator, username string) {
	v.CheckStrNotEmpty(username, "username")
	v.Check(len(username) <= 30, "username", "must not be more than 30 characters long")
}

func ValidateUserEmail(v *validator.Validator, email string) {
	v.CheckStrNotEmpty(email, "email")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUserPasswordPlaintext(v *validator.Validator, pt string) {
	v.CheckStrNotEmpty(pt, "password")
	v.Check(len(pt) >= 8, "password", "must be at least 8 characters long")
	v.Check(len(pt) <= 72, "password", "must not be more than 72 characters long")
}

func ValidateName(v *validator.Validator, name string) {
	v.CheckStrNotEmpty(name, "name")
	v.Check(len(name) <= 30, "name", "must not be more than 30 characters long")
}

func ExtractUserValidationErrors(v *validator.Validator, err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		if pgErr.ConstraintName == UserConstraint.UniqueEmail {
			v.AddError("email", "duplicate email")
		}
		if pgErr.ConstraintName == UserConstraint.UniqueUsername {
			v.AddError("username", "duplicate username")
		}
	}
}

func (u UserInput) Validate(v *validator.Validator) {
	ValidateUsername(v, u.Username)
	ValidateName(v, u.Name)
	ValidateUserEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidateUserPasswordPlaintext(v, *u.Password.plaintext)
		if !v.Valid() {
			return
		}
	}

	if u.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
