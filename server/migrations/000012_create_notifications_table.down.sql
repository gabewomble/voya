DROP TABLE IF EXISTS notifications;

DROP TYPE IF EXISTS notification_type;

-- Rename column if it exists
ALTER TABLE
    trip_members RENAME COLUMN updated_by TO removed_by;

ALTER TABLE
    trip_members
ADD
    COLUMN IF NOT EXISTS removed_at TIMESTAMP WITH TIME ZONE;