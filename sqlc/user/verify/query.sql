-- name: Get :one

SELECT *
FROM "user".verification
WHERE id = $1;

-- name: Create :one

INSERT INTO "user".verification (user_id)
VALUES ($1) RETURNING *;
