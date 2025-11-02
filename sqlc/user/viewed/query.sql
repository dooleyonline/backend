-- name: GetAll :many
SELECT
  *
FROM
  "user"."viewed";

-- name: Create :one
INSERT INTO
"user"."viewed" (user_id, item_id, created_at)
VALUES
($1, $2, $3)
RETURNING *;

-- name: Delete :one
DELETE FROM
"user"."viewed"
WHERE
user_id = $1
RETURNING *;
