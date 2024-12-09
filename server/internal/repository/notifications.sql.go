// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: notifications.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countUnreadNotifications = `-- name: CountUnreadNotifications :one
SELECT
    COUNT(*)
FROM
    notifications
WHERE
    user_id = $1
    AND read_at IS NULL
`

func (q *Queries) CountUnreadNotifications(ctx context.Context, userID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countUnreadNotifications, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteNotification = `-- name: DeleteNotification :exec
DELETE FROM
    notifications
WHERE
    id = $1
    AND user_id = $2
`

type DeleteNotificationParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteNotification(ctx context.Context, arg DeleteNotificationParams) error {
	_, err := q.db.Exec(ctx, deleteNotification, arg.ID, arg.UserID)
	return err
}

const deleteNotificationsByType = `-- name: DeleteNotificationsByType :exec
DELETE FROM
    notifications
WHERE
    user_id = $1
    AND trip_id = $2
    AND notification_type = $3
`

type DeleteNotificationsByTypeParams struct {
	UserID uuid.UUID        `json:"user_id"`
	TripID uuid.UUID        `json:"trip_id"`
	Type   NotificationType `json:"type"`
}

func (q *Queries) DeleteNotificationsByType(ctx context.Context, arg DeleteNotificationsByTypeParams) error {
	_, err := q.db.Exec(ctx, deleteNotificationsByType, arg.UserID, arg.TripID, arg.Type)
	return err
}

const getNotificationById = `-- name: GetNotificationById :one
SELECT
    id,
    user_id,
    trip_id,
    message,
    notification_type,
    created_at,
    read_at,
    target_user_id,
    created_by
FROM
    notifications
WHERE
    id = $1
    AND user_id = $2
`

type GetNotificationByIdParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type GetNotificationByIdRow struct {
	ID               uuid.UUID        `json:"id"`
	UserID           uuid.UUID        `json:"user_id"`
	TripID           uuid.UUID        `json:"trip_id"`
	Message          string           `json:"message"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt        *time.Time       `json:"created_at"`
	ReadAt           *time.Time       `json:"read_at"`
	TargetUserID     uuid.UUID        `json:"target_user_id"`
	CreatedBy        uuid.UUID        `json:"created_by"`
}

func (q *Queries) GetNotificationById(ctx context.Context, arg GetNotificationByIdParams) (GetNotificationByIdRow, error) {
	row := q.db.QueryRow(ctx, getNotificationById, arg.ID, arg.UserID)
	var i GetNotificationByIdRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TripID,
		&i.Message,
		&i.NotificationType,
		&i.CreatedAt,
		&i.ReadAt,
		&i.TargetUserID,
		&i.CreatedBy,
	)
	return i, err
}

const getUnreadNotifications = `-- name: GetUnreadNotifications :many
SELECT
    id,
    user_id,
    trip_id,
    message,
    notification_type,
    created_at,
    read_at,
    target_user_id,
    created_by
FROM
    notifications
WHERE
    user_id = $1
    AND read_at IS NULL
ORDER BY
    created_at DESC
`

type GetUnreadNotificationsRow struct {
	ID               uuid.UUID        `json:"id"`
	UserID           uuid.UUID        `json:"user_id"`
	TripID           uuid.UUID        `json:"trip_id"`
	Message          string           `json:"message"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt        *time.Time       `json:"created_at"`
	ReadAt           *time.Time       `json:"read_at"`
	TargetUserID     uuid.UUID        `json:"target_user_id"`
	CreatedBy        uuid.UUID        `json:"created_by"`
}

func (q *Queries) GetUnreadNotifications(ctx context.Context, userID uuid.UUID) ([]GetUnreadNotificationsRow, error) {
	rows, err := q.db.Query(ctx, getUnreadNotifications, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUnreadNotificationsRow
	for rows.Next() {
		var i GetUnreadNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TripID,
			&i.Message,
			&i.NotificationType,
			&i.CreatedAt,
			&i.ReadAt,
			&i.TargetUserID,
			&i.CreatedBy,
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

const insertNotification = `-- name: InsertNotification :exec
INSERT INTO
    notifications (
        user_id,
        trip_id,
        message,
        notification_type,
        target_user_id,
        created_by
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
`

type InsertNotificationParams struct {
	UserID       uuid.UUID        `json:"user_id"`
	TripID       uuid.UUID        `json:"trip_id"`
	Message      string           `json:"message"`
	Type         NotificationType `json:"type"`
	TargetUserID uuid.UUID        `json:"target_user_id"`
	CreatedBy    uuid.UUID        `json:"created_by"`
}

func (q *Queries) InsertNotification(ctx context.Context, arg InsertNotificationParams) error {
	_, err := q.db.Exec(ctx, insertNotification,
		arg.UserID,
		arg.TripID,
		arg.Message,
		arg.Type,
		arg.TargetUserID,
		arg.CreatedBy,
	)
	return err
}

const listNotifications = `-- name: ListNotifications :many
SELECT
    id,
    user_id,
    trip_id,
    message,
    notification_type,
    created_at,
    read_at,
    target_user_id,
    created_by
FROM
    notifications
WHERE
    user_id = $1
ORDER BY
    created_at DESC
LIMIT
    $3 OFFSET $2
`

type ListNotificationsParams struct {
	UserID             uuid.UUID `json:"user_id"`
	NotificationOffset int32     `json:"notification_offset"`
	NotificationLimit  int32     `json:"notification_limit"`
}

type ListNotificationsRow struct {
	ID               uuid.UUID        `json:"id"`
	UserID           uuid.UUID        `json:"user_id"`
	TripID           uuid.UUID        `json:"trip_id"`
	Message          string           `json:"message"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt        *time.Time       `json:"created_at"`
	ReadAt           *time.Time       `json:"read_at"`
	TargetUserID     uuid.UUID        `json:"target_user_id"`
	CreatedBy        uuid.UUID        `json:"created_by"`
}

func (q *Queries) ListNotifications(ctx context.Context, arg ListNotificationsParams) ([]ListNotificationsRow, error) {
	rows, err := q.db.Query(ctx, listNotifications, arg.UserID, arg.NotificationOffset, arg.NotificationLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNotificationsRow
	for rows.Next() {
		var i ListNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TripID,
			&i.Message,
			&i.NotificationType,
			&i.CreatedAt,
			&i.ReadAt,
			&i.TargetUserID,
			&i.CreatedBy,
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

const markNotificationAsRead = `-- name: MarkNotificationAsRead :exec
UPDATE
    notifications
SET
    read_at = NOW()
WHERE
    id = $1
    AND user_id = $2
`

type MarkNotificationAsReadParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) MarkNotificationAsRead(ctx context.Context, arg MarkNotificationAsReadParams) error {
	_, err := q.db.Exec(ctx, markNotificationAsRead, arg.ID, arg.UserID)
	return err
}

const markNotificationsAsRead = `-- name: MarkNotificationsAsRead :exec
UPDATE
    notifications
SET
    read_at = NOW()
WHERE
    user_id = $1
    AND read_at IS NULL
`

func (q *Queries) MarkNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, markNotificationsAsRead, userID)
	return err
}

const notifyTripMembers = `-- name: NotifyTripMembers :exec
INSERT INTO
    notifications (
        user_id,
        trip_id,
        message,
        notification_type,
        target_user_id,
        created_by
    )
SELECT
    tm.user_id,
    $1,
    $2,
    $3,
    $4,
    $5
FROM
    trip_members tm
WHERE
    tm.trip_id = $1
    AND tm.member_status IN ('accepted', 'owner')
`

type NotifyTripMembersParams struct {
	TripID       uuid.UUID        `json:"trip_id"`
	Message      string           `json:"message"`
	Type         NotificationType `json:"type"`
	TargetUserID uuid.UUID        `json:"target_user_id"`
	CreatedBy    uuid.UUID        `json:"created_by"`
}

func (q *Queries) NotifyTripMembers(ctx context.Context, arg NotifyTripMembersParams) error {
	_, err := q.db.Exec(ctx, notifyTripMembers,
		arg.TripID,
		arg.Message,
		arg.Type,
		arg.TargetUserID,
		arg.CreatedBy,
	)
	return err
}
