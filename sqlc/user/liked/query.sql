-- name: GetAll :many
SELECT
  *
FROM
  "user"."liked";

-- name: Like :exec
INSERT INTO
"user"."liked" (user_id, item_id)
VALUES
($1, $2)
RETURNING *;

-- name: Unlike :exec
DELETE FROM
"user"."liked"
WHERE
user_id = $1 AND item_id = $2
RETURNING *;
