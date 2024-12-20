-- name: InsertTrip :one
INSERT INTO
    trips (name, description)
VALUES
    (@name, @description) RETURNING *;

-- name: GetTripById :one
SELECT
    *
FROM
    trips
WHERE
    id = @id
    AND (
        id IN (
            SELECT
                trip_id
            FROM
                trip_members
            WHERE
                user_id = @user_id
        )
    );

-- name: DeleteTripById :exec
DELETE FROM
    trips
WHERE
    id = @id
    AND id IN (
        SELECT
            trip_id
        FROM
            trip_members
        WHERE
            user_id = @user_id
            AND member_status = 'owner'
    );

-- name: ListTrips :many
SELECT
    *
FROM
    trips
WHERE
    id IN (
        SELECT
            trip_id
        FROM
            trip_members
        WHERE
            user_id = @user_id
            AND member_status IN ('accepted', 'owner')
    );

-- name: CheckUserCanViewTrip :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            trips
        WHERE
            id = @id
            AND id IN (
                SELECT
                    trip_id
                FROM
                    trip_members
                WHERE
                    user_id = @user_id
            )
    );

-- name: CheckUserCanEditTrip :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            trips
        WHERE
            id = @id
            AND id IN (
                SELECT
                    trip_id
                FROM
                    trip_members
                WHERE
                    user_id = @user_id
                    AND (
                        member_status = 'accepted'
                        OR member_status = 'owner'
                    )
            )
    );