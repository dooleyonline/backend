-- name: GetAll :many
SELECT
  *
FROM
  "user"."viewed";

-- name: Create :exec
INSERT INTO
"user"."viewed" (user_id, item_id)
VALUES
($1, $2)
RETURNING *;

-- name: GetByUserID :many
SELECT
  *
FROM
  "user"."viewed"
WHERE
  user_id = $1;
