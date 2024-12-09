ALTER TABLE
    notifications DROP COLUMN created_by;

ALTER TABLE
    notifications DROP COLUMN target_user_id;

ALTER TABLE
    notifications
ADD
    COLUMN metadata JSONB;