-- name: InsertUser :one
INSERT INTO
    users (name, email, username, password_hash, activated)
VALUES
    (
        @name,
        @email,
        @username,
        @password_hash,
        @activated
    ) RETURNING id,
    created_at,
    version;

-- name: GetUserByEmail :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    email = @email;

-- name: GetUserByUsername :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    username = @username;

-- name: GetUserById :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    id = @id;

-- name: GetUserForToken :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
    INNER JOIN tokens ON users.id = tokens.user_id
WHERE
    tokens.hash = @token_hash
    AND tokens.scope = @token_scope
    AND tokens.expiry > @token_expiry;

-- name: GetUserForRefreshToken :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
    INNER JOIN tokens ON users.id = tokens.user_id
WHERE
    tokens.refresh_token = @refresh_token
    AND tokens.scope = @token_scope
    AND tokens.expiry > @token_expiry;

-- name: UpdateUser :one
UPDATE
    users
SET
    name = @name,
    email = @email,
    username = @username,
    password_hash = @password_hash,
    activated = @activated,
    version = version + 1
WHERE
    id = @id
    AND version = @version RETURNING version;

-- name: CheckUsernameExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            username = @username
    );

-- name: SearchUsers :many
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version,
    username
FROM
    users
WHERE
    (
        name ILIKE '%' || @identifier || '%'
        OR email ILIKE '%' || @identifier || '%'
        OR username ILIKE '%' || @identifier || '%'
    )
    AND (id != @user_id)
LIMIT
    @user_limit;

-- name: CheckUserExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            id = @id
    );

-- name: SearchUsersNotInTrip :many
SELECT
    u.id,
    u.created_at,
    u.name,
    u.email,
    u.password_hash,
    u.activated,
    u.version,
    u.username
FROM
    users u
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            trip_members tm
        WHERE
            tm.trip_id = @trip_id
            AND tm.user_id = u.id
            AND tm.member_status IN ('accepted', 'owner', 'pending')
    )
    AND (
        u.name ILIKE '%' || @identifier || '%'
        OR u.email ILIKE '%' || @identifier || '%'
        OR u.username ILIKE '%' || @identifier || '%'
    )
LIMIT
    @user_limit;