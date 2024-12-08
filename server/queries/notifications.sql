-- name: InsertNotification :exec
INSERT INTO
    notifications (user_id, trip_id, message, TYPE, metadata)
VALUES
    (@user_id, @trip_id, @message, @type, @metadata);

-- name: GetUnreadNotifications :many
SELECT
    id,
    user_id,
    trip_id,
    message,
    TYPE,
    created_at,
    read_at,
    metadata
FROM
    notifications
WHERE
    user_id = @user_id
    AND read_at IS NULL
ORDER BY
    created_at DESC;

-- name: MarkNotificationAsRead :exec
UPDATE
    notifications
SET
    read_at = NOW()
WHERE
    id = @id
    AND user_id = @user_id;

-- name: GetNotificationById :one
SELECT
    id,
    user_id,
    trip_id,
    message,
    TYPE,
    created_at,
    read_at,
    metadata
FROM
    notifications
WHERE
    id = @id
    AND user_id = @user_id;

-- name: DeleteNotification :exec
DELETE FROM
    notifications
WHERE
    id = @id
    AND user_id = @user_id;

-- name: ListNotifications :many
SELECT
    id,
    user_id,
    trip_id,
    message,
    TYPE,
    created_at,
    read_at,
    metadata
FROM
    notifications
WHERE
    user_id = @user_id
ORDER BY
    created_at DESC
LIMIT
    @notification_limit OFFSET @notification_offset;

-- name: CountUnreadNotifications :one
SELECT
    COUNT(*)
FROM
    notifications
WHERE
    user_id = @user_id
    AND read_at IS NULL;