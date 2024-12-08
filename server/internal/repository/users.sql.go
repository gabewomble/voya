// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkUserExists = `-- name: CheckUserExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            id = $1
    )
`

func (q *Queries) CheckUserExists(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, checkUserExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUsernameExists = `-- name: CheckUsernameExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            username = $1
    )
`

func (q *Queries) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRow(ctx, checkUsernameExists, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.Username,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.Username,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.Username,
	)
	return i, err
}

const getUserForRefreshToken = `-- name: GetUserForRefreshToken :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
    INNER JOIN tokens ON users.id = tokens.user_id
WHERE
    tokens.refresh_token = $1
    AND tokens.scope = $2
    AND tokens.expiry > $3
`

type GetUserForRefreshTokenParams struct {
	RefreshToken []byte    `json:"refresh_token"`
	TokenScope   string    `json:"token_scope"`
	TokenExpiry  time.Time `json:"token_expiry"`
}

func (q *Queries) GetUserForRefreshToken(ctx context.Context, arg GetUserForRefreshTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserForRefreshToken, arg.RefreshToken, arg.TokenScope, arg.TokenExpiry)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.Username,
	)
	return i, err
}

const getUserForToken = `-- name: GetUserForToken :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
    INNER JOIN tokens ON users.id = tokens.user_id
WHERE
    tokens.hash = $1
    AND tokens.scope = $2
    AND tokens.expiry > $3
`

type GetUserForTokenParams struct {
	TokenHash   []byte    `json:"token_hash"`
	TokenScope  string    `json:"token_scope"`
	TokenExpiry time.Time `json:"token_expiry"`
}

func (q *Queries) GetUserForToken(ctx context.Context, arg GetUserForTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserForToken, arg.TokenHash, arg.TokenScope, arg.TokenExpiry)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.Username,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO
    users (name, email, username, password_hash, activated)
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING id,
    created_at,
    version
`

type InsertUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
}

type InsertUserRow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (InsertUserRow, error) {
	row := q.db.QueryRow(ctx, insertUser,
		arg.Name,
		arg.Email,
		arg.Username,
		arg.PasswordHash,
		arg.Activated,
	)
	var i InsertUserRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Version)
	return i, err
}

const searchUsers = `-- name: SearchUsers :many
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    (
        name ILIKE '%' || $1 || '%'
        OR email ILIKE '%' || $1 || '%'
        OR username ILIKE '%' || $1 || '%'
    )
    AND (id != $2)
LIMIT
    $3
`

type SearchUsersParams struct {
	Identifier pgtype.Text `json:"identifier"`
	UserID     uuid.UUID   `json:"user_id"`
	UserLimit  int32       `json:"user_limit"`
}

func (q *Queries) SearchUsers(ctx context.Context, arg SearchUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, searchUsers, arg.Identifier, arg.UserID, arg.UserLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Email,
			&i.PasswordHash,
			&i.Activated,
			&i.Version,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE
    users
SET
    name = $1,
    email = $2,
    username = $3,
    password_hash = $4,
    activated = $5,
    version = version + 1
WHERE
    id = $6
    AND version = $7 RETURNING version
`

type UpdateUserParams struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"password_hash"`
	Activated    bool      `json:"activated"`
	ID           uuid.UUID `json:"id"`
	Version      int32     `json:"version"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (int32, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.Username,
		arg.PasswordHash,
		arg.Activated,
		arg.ID,
		arg.Version,
	)
	var version int32
	err := row.Scan(&version)
	return version, err
}
