// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: trips.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const addTripMember = `-- name: AddTripMember :exec
INSERT INTO
    trip_members (trip_id, user_id)
VALUES
    ($1, $2)
`

type AddTripMemberParams struct {
	TripID uuid.UUID `json:"trip_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) AddTripMember(ctx context.Context, arg AddTripMemberParams) error {
	_, err := q.db.Exec(ctx, addTripMember, arg.TripID, arg.UserID)
	return err
}

const deleteTripById = `-- name: DeleteTripById :exec
DELETE FROM
    trips
WHERE
    id = $1
    AND owner_id = $2
`

type DeleteTripByIdParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteTripById(ctx context.Context, arg DeleteTripByIdParams) error {
	_, err := q.db.Exec(ctx, deleteTripById, arg.ID, arg.UserID)
	return err
}

const getTripById = `-- name: GetTripById :one
SELECT
    id, name, description, start_date, end_date, created_at, updated_at, owner_id
FROM
    trips
WHERE
    id = $1
    AND (
        owner_id = $2
        OR id IN (
            SELECT
                trip_id
            FROM
                trip_members
            WHERE
                user_id = $2
        )
    )
`

type GetTripByIdParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetTripById(ctx context.Context, arg GetTripByIdParams) (Trip, error) {
	row := q.db.QueryRow(ctx, getTripById, arg.ID, arg.UserID)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OwnerID,
	)
	return i, err
}

const getTripMembers = `-- name: GetTripMembers :many
SELECT
    u.id,
    u.name,
    u.email
FROM
    users u
    INNER JOIN trip_members tm ON u.id = tm.user_id
WHERE
    tm.trip_id = $1
`

type GetTripMembersRow struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (q *Queries) GetTripMembers(ctx context.Context, tripID uuid.UUID) ([]GetTripMembersRow, error) {
	rows, err := q.db.Query(ctx, getTripMembers, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTripMembersRow
	for rows.Next() {
		var i GetTripMembersRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertTrip = `-- name: InsertTrip :one
INSERT INTO
    trips (name, description, owner_id)
VALUES
    ($1, $2, $3) RETURNING id, name, description, start_date, end_date, created_at, updated_at, owner_id
`

type InsertTripParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	OwnerID     uuid.UUID   `json:"owner_id"`
}

func (q *Queries) InsertTrip(ctx context.Context, arg InsertTripParams) (Trip, error) {
	row := q.db.QueryRow(ctx, insertTrip, arg.Name, arg.Description, arg.OwnerID)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OwnerID,
	)
	return i, err
}

const listTrips = `-- name: ListTrips :many
SELECT
    id, name, description, start_date, end_date, created_at, updated_at, owner_id
FROM
    trips
WHERE
    owner_id = $1
    OR id IN (
        SELECT
            trip_id
        FROM
            trip_members
        WHERE
            user_id = $1
    )
`

func (q *Queries) ListTrips(ctx context.Context, userID uuid.UUID) ([]Trip, error) {
	rows, err := q.db.Query(ctx, listTrips, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Trip
	for rows.Next() {
		var i Trip
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OwnerID,
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

const removeTripMember = `-- name: RemoveTripMember :exec
DELETE FROM
    trip_members
WHERE
    trip_id = $1
    AND user_id = $2
`

type RemoveTripMemberParams struct {
	TripID uuid.UUID `json:"trip_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) RemoveTripMember(ctx context.Context, arg RemoveTripMemberParams) error {
	_, err := q.db.Exec(ctx, removeTripMember, arg.TripID, arg.UserID)
	return err
}
