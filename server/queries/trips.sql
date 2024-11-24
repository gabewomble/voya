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
        owner_id = @user_id
        OR id IN (
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
    owner_id = @user_id
    OR id IN (
        SELECT
            trip_id
        FROM
            trip_members
        WHERE
            user_id = @user_id
    );

-- name: AddTripMember :exec
INSERT INTO
    trip_members (trip_id, user_id)
VALUES
    (@trip_id, @user_id);

-- name: RemoveTripMember :exec
DELETE FROM
    trip_members
WHERE
    trip_id = @trip_id
    AND user_id = @user_id;

-- name: GetTripMembers :many
SELECT
    u.id,
    u.name,
    u.email
FROM
    users u
    INNER JOIN trip_members tm ON u.id = tm.user_id
WHERE
    tm.trip_id = @trip_id;