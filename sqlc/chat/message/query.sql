-- name: Create :exec
INSERT INTO
  chat.message(id, room_id, sent_by, body)
VALUES
  ($1, $2, $3, $4);

-- name: GetMany :many
SELECT
  *
FROM
  chat.message
WHERE
  room_id = $1
ORDER BY
  sent_at DESC
LIMIT $2;

-- name: Delete :exec
DELETE FROM
  chat.message
WHERE
  id = $1;

-- name: EditMessage :exec
UPDATE
  chat.message
SET
  body = $2, edited = true
WHERE
  id = $1;