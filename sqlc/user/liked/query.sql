-- name: GetAll :many
SELECT
  *
FROM
  "user"."liked";

-- name: Create :one
INSERT INTO
"user"."liked" (user_id, item_id, created_at)
VALUES
($1, $2, $3)
RETURNING *;

-- name: Delete :one
DELETE FROM
"user"."liked"
WHERE
user_id = $1
RETURNING *;
