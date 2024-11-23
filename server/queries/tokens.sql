-- name: InsertToken :exec
INSERT INTO
    tokens (hash, user_id, expiry, scope, refresh_token)
VALUES
    (
        @token_hash,
        @user_id,
        @token_expiry,
        @token_scope,
        @refresh_token
    );

-- name: DeleteAllTokensForUser :exec
DELETE FROM
    tokens
WHERE
    scope = @token_scope
    AND user_id = @user_id;

-- name: DeleteToken :exec
DELETE FROM
    tokens
WHERE
    hash = @token_hash;

-- name: DeleteExpiredTokensForUser :exec
DELETE FROM
    tokens
WHERE
    user_id = @user_id
    AND expiry < NOW();