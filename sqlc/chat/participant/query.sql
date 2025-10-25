-- name: Create :exec
INSERT INTO
  chat.participant (room_id, user_id)
VALUES ($1, $2);


-- name: Get :one
SELECT
  *
FROM
  chat.participant
WHERE
  room_id = $1
  AND user_id = $2; 


-- name: Delete :exec
DELETE FROM
  chat.participant
WHERE
  room_id = $1
  AND user_id = $2;


-- name: GetByUserID :many
SELECT 
  *
FROM
  chat.participant
WHERE
  user_id = $1;

-- name: GetByRoomID :many
SELECT 
  *
FROM
  chat.participant
WHERE
  room_id = $1; 