-- name: InsertNotification :exec
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
        @user_id,
        @trip_id,
        @message,
        @type,
        @target_user_id,
        @created_by
    );

-- name: NotifyOtherTripMembers :exec
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
    @trip_id,
    @message,
    @type,
    @target_user_id,
    @created_by
FROM
    trip_members tm
WHERE
    tm.trip_id = @trip_id
    AND tm.member_status IN ('accepted', 'owner')
    AND tm.user_id != @created_by;

-- name: GetUnreadNotifications :many
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
    target_user_id,
    created_by
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
    target_user_id,
    created_by
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
