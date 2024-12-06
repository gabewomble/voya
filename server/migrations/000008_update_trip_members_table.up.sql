ALTER TABLE
    trip_members
ADD
    COLUMN invited_by UUID REFERENCES users(id),
ADD
    COLUMN member_status TEXT NOT NULL DEFAULT 'pending',
ADD
    COLUMN removed_by UUID REFERENCES users(id),
ADD
    COLUMN removed_at TIMESTAMP WITH TIME ZONE;