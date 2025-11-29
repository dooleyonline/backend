-- name: Create :exec
INSERT INTO
"user"."viewed" (user_id, item_id)
VALUES
($1, $2)
RETURNING *;

-- name: GetViewed :many
SELECT
item_id
FROM "user"."viewed"
WHERE user_id = $1;
