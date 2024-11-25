ALTER TABLE
    users
ALTER COLUMN
    username DROP NOT NULL;

ALTER TABLE
    users
ALTER COLUMN
    username
SET
    DEFAULT 'default_username';

ALTER TABLE
    users DROP CONSTRAINT unique_username;