ALTER TABLE
    notifications
ADD
    COLUMN created_by UUID REFERENCES users(id);

ALTER TABLE
    notifications
ADD
    COLUMN target_user_id UUID REFERENCES users(id);

ALTER TABLE
    notifications DROP COLUMN metadata;