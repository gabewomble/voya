-- 1. Remove the owner_in_trip_members constraint
ALTER TABLE
    trips DROP CONSTRAINT IF EXISTS owner_in_trip_members;

-- 2. Convert the member_status column back to VARCHAR
ALTER TABLE
    trip_members
ALTER COLUMN
    member_status TYPE TEXT NOT NULL DEFAULT 'pending';

-- 3. Drop the member_status_enum type
DROP TYPE IF EXISTS member_status_enum;