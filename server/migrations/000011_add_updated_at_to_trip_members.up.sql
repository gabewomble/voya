ALTER TABLE
    trip_members
ADD
    COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;