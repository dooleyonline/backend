-- name: GetBlocksByBlockerID :many
SELECT
  blocked_id
FROM
  "user"."block"
WHERE blocker_id=$1;

-- name: GetBlocksByBlockedID :many
SELECT
  blocker_id
FROM
  "user"."block"
WHERE blocked_id=$1;

-- name: Block :exec
INSERT INTO
"user"."block" (blocker_id, blocked_id)
VALUES
($1, $2)
RETURNING *;

-- name: Unblock :exec
DELETE FROM
"user"."block"
WHERE
blocker_id = $1 AND blocked_id = $2
RETURNING *;
