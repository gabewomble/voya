-- 1. Backfill missing trip_members entries for owners
INSERT INTO
    trip_members (user_id, trip_id, member_status)
SELECT
    owner_id AS user_id,
    id AS trip_id,
    'owner' AS member_status
FROM
    trips
WHERE
    owner_id IS NOT NULL
    AND NOT EXISTS (
        SELECT
            1
        FROM
            trip_members tm
        WHERE
            tm.user_id = trips.owner_id
            AND tm.trip_id = trips.id
    );

-- 2. Create the member_status_enum type
CREATE TYPE member_status_enum AS ENUM (
    'owner',
    'pending',
    'accepted',
    'declined',
    'removed',
    'cancelled'
);

-- 3. Convert the member_status column to use the new enum
ALTER TABLE
    trip_members
ALTER COLUMN
    member_status DROP DEFAULT;

ALTER TABLE
    trip_members
ALTER COLUMN
    member_status TYPE member_status_enum USING member_status :: member_status_enum;

ALTER TABLE
    trip_members
ALTER COLUMN
    member_status
SET
    DEFAULT 'pending';