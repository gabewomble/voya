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

const checkUserCanEditTrip = `-- name: CheckUserCanEditTrip :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            trips
        WHERE
            id = $1
            AND id IN (
                SELECT
                    trip_id
                FROM
                    trip_members
                WHERE
                    user_id = $2
                    AND (
                        member_status = 'accepted'
                        OR member_status = 'owner'
                    )
            )
    )
`

type CheckUserCanEditTripParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CheckUserCanEditTrip(ctx context.Context, arg CheckUserCanEditTripParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkUserCanEditTrip, arg.ID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUserCanViewTrip = `-- name: CheckUserCanViewTrip :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            trips
        WHERE
            id = $1
            AND id IN (
                SELECT
                    trip_id
                FROM
                    trip_members
                WHERE
                    user_id = $2
            )
    )
`

type CheckUserCanViewTripParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CheckUserCanViewTrip(ctx context.Context, arg CheckUserCanViewTripParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkUserCanViewTrip, arg.ID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const deleteTripById = `-- name: DeleteTripById :exec
DELETE FROM
    trips
WHERE
    id = $1
    AND id IN (
        SELECT
            trip_id
        FROM
            trip_members
        WHERE
            user_id = $2
            AND member_status = 'owner'
    )
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
    id, name, description, start_date, end_date, created_at, updated_at
FROM
    trips
WHERE
    id = $1
    AND (
        id IN (
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
	)
	return i, err
}

const insertTrip = `-- name: InsertTrip :one
INSERT INTO
    trips (name, description)
VALUES
    ($1, $2) RETURNING id, name, description, start_date, end_date, created_at, updated_at
`

type InsertTripParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) InsertTrip(ctx context.Context, arg InsertTripParams) (Trip, error) {
	row := q.db.QueryRow(ctx, insertTrip, arg.Name, arg.Description)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTrips = `-- name: ListTrips :many
SELECT
    id, name, description, start_date, end_date, created_at, updated_at
FROM
    trips
WHERE
    id IN (
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
