-- name: CreateRoom :one
INSERT INTO
  "chat"."room" (participants)
VALUES
  ($1)
RETURNING id;

-- name: AddParticipant :exec
UPDATE
  chat.room
SET
  participants = $1
WHERE
  id = $2;

-- name: RemoveParticipant :exec
UPDATE
  chat.room
SET
  participants = $1
WHERE
  id = $2;

-- name: AddMessage :exec
UPDATE
  chat.room
SET
  messages = array_append(messages, @message_id::bigint)
WHERE
  id = $1
  AND NOT (messages @> ARRAY[@message_id::bigint]);
