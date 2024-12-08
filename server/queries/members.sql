-- name: InsertTripOwner :exec
INSERT INTO
    trip_members (trip_id, user_id, invited_by, member_status)
VALUES
    (@trip_id, @owner_id, NULL, 'owner');

-- name: AddUserToTrip :exec
INSERT INTO
    trip_members (trip_id, user_id, invited_by, member_status)
VALUES
    (@trip_id, @user_id, @invited_by, 'pending') ON CONFLICT (trip_id, user_id) DO
UPDATE
SET
    invited_by = EXCLUDED.invited_by,
    member_status = 'pending',
    removed_by = NULL,
    removed_at = NULL;

-- name: UpdateTripMemberStatus :exec
UPDATE
    trip_members
SET
    member_status = @member_status,
    removed_by = @removed_by,
    removed_at = @removed_at,
    updated_at = CURRENT_TIMESTAMP
WHERE
    trip_id = @trip_id
    AND user_id = @user_id;

-- name: GetTripMembers :many
SELECT
    u.id,
    u.name,
    u.email,
    tm.updated_at,
    tm.member_status
FROM
    users u
    INNER JOIN trip_members tm ON u.id = tm.user_id
WHERE
    tm.trip_id = @trip_id;