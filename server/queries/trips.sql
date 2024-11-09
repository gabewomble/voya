-- name: InsertTrip :one
INSERT INTO
    trips (name, description)
VALUES
    ($1, $2) RETURNING *;

-- name: GetTripById :one
SELECT
    *
FROM
    trips
WHERE
    id = @id;

-- name: DeleteTripById :exec
DELETE FROM
    trips
WHERE
    id = @id;

-- name: ListTrips :many
SELECT
    *
FROM
    trips;