-- name: InsertToken :exec
INSERT INTO
    tokens (hash, user_id, expiry, scope)
VALUES
    ($1, $2, $3, $4);

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
