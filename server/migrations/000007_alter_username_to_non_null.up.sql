ALTER TABLE
    users
ALTER COLUMN
    username
SET
    NOT NULL;

ALTER TABLE
    users
ALTER COLUMN
    username DROP DEFAULT;

ALTER TABLE
    users
ADD
    CONSTRAINT username_unique UNIQUE (username);