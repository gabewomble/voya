ALTER TABLE
    trips
ADD
    COLUMN owner_id UUID REFERENCES users(id) ON DELETE CASCADE;