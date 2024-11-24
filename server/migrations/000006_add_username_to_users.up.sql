ALTER TABLE
    users
ADD
    COLUMN username citext DEFAULT 'default_username' NOT NULL;

CREATE UNIQUE INDEX idx_users_username ON users (username);