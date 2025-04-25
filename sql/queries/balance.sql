-- name: GetBalanceById :one
SELECT * FROM balance WHERE id = $1;

-- name: CreateBalance :one
INSERT INTO balance (id, user_id, balance) VALUES ($1, $2, $3) RETURNING *;

