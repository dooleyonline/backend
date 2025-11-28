-- name: CreateRoom :one
INSERT INTO
  chat.room (participants)
VALUES
  ($1)
RETURNING id;

-- name: DeleteRoom :exec 
DELETE FROM
  chat.room
WHERE
  id = $1;

-- name: GetRoomByID :one
SELECT
  *
FROM
  chat.room
WHERE
  id = $1;

-- name: AddParticipant :exec
UPDATE
  chat.room
SET
  participants = array_append(participants, @user_id::uuid)
WHERE
  id = @room_id
  AND NOT participants @> ARRAY[@user_id::uuid];

-- name: RemoveParticipant :exec
UPDATE
  chat.room
SET
  participants = array_remove(participants, @user_id::uuid)
WHERE
  id = @room_id
  AND participants @> ARRAY[@user_id::uuid];

