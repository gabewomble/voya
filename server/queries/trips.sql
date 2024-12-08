-- name: InsertTrip :one
INSERT INTO
    trips (name, description, owner_id)
VALUES
    (@name, @description, @owner_id) RETURNING *;

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
    AND owner_id = @user_id;

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
    );