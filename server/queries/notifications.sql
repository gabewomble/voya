-- name: InsertNotification :exec
INSERT INTO
    notifications (
        user_id,
        trip_id,
        message,
        notification_type,
        metadata
    )
VALUES
    (@user_id, @trip_id, @message, @type, @metadata);

-- name: GetUnreadNotifications :many
SELECT
    id,
    user_id,
    trip_id,
    message,
    notification_type,
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

-- name: MarkNotificationsAsRead :exec
UPDATE
    notifications
SET
    read_at = NOW()
WHERE
    user_id = @user_id
    AND read_at IS NULL;

-- name: GetNotificationById :one
SELECT
    id,
    user_id,
    trip_id,
    message,
    notification_type,
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
    notification_type,
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

-- name: DeleteNotificationsByType :exec
DELETE FROM
    notifications
WHERE
    user_id = @user_id
    AND trip_id = @trip_id
    AND notification_type = @type;