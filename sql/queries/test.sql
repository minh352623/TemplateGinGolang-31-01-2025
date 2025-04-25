-- name: GetTestById :one
SELECT * FROM test WHERE id = $1;

-- name: CreateTest :one
INSERT INTO test (id, name, user_id, balance) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateTest :one
UPDATE test SET name = $2, balance = $3 WHERE id = $1 RETURNING *;

-- name: DeleteTest :one
DELETE FROM test WHERE id = $1 RETURNING *;

-- name: GetAllTestPaginationAndSearch :many
SELECT * FROM test
WHERE name ILIKE '%' || $1 || '%' AND id = $2
ORDER BY id
LIMIT $3 OFFSET $4;
