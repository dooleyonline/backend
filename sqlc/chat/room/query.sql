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


-- name: IncrementMessageCount :exec
UPDATE
  chat.room
SET
  message_count = message_count + 1
WHERE
  id = $1;

-- name: SyncAllMessageCounts :exec
UPDATE
  chat.room room
SET
  message_count = COALESCE(m.count, 0)
FROM (
  SELECT room_id, COUNT(*) as count
  FROM chat.message
  GROUP BY room_id
) m
WHERE room.id = m.room_id; 