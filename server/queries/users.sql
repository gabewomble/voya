-- name: InsertUser :one
INSERT INTO
    users (name, email, password_hash, activated)
VALUES
    (@name, @email, @password_hash, @activated) RETURNING id,
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
    version
FROM
    users
WHERE
    email = @email;

-- name: GetUserById :one
SELECT
    id,
    created_at,
    name,
    email,
    password_hash,
    activated,
    version
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
    version
FROM
    users
    INNER JOIN tokens ON users.id = tokens.user_id
WHERE
    tokens.hash = @token_hash
    AND tokens.scope = @token_scope
    AND tokens.expiry > @token_expiry;

-- name: UpdateUser :one
UPDATE
    users
SET
    name = @name,
    email = @email,
    password_hash = @password_hash,
    activated = @activated,
    version = version + 1
WHERE
    id = @id
    AND version = @version RETURNING version;