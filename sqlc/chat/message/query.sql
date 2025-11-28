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
  room_id = @room_id
ORDER BY
  sent_at DESC
LIMIT 
  10 OFFSET 10*(CAST(@page AS integer)-1);


-- name: GetLatestMessage :one
SELECT 
  *
FROM
  chat.message
WHERE
  room_id = $1
ORDER BY 
  sent_at DESC
LIMIT 1;

-- name: GetByID :one
SELECT
  *
FROM
  chat.message
WHERE
  id = $1;


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

