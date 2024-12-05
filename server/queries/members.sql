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
    removed_at = @removed_at
WHERE
    trip_id = @trip_id
    AND user_id = @user_id;