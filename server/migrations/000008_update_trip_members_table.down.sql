ALTER TABLE
    trip_members DROP COLUMN IF EXISTS invited_by,
    DROP COLUMN IF EXISTS member_status,
    DROP COLUMN IF EXISTS removed_by,
    DROP COLUMN IF EXISTS removed_at;