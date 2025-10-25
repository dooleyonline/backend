-- name: Create :exec
INSERT INTO
  chat.participant (room_id, user_id)
VALUES ($1, $2);
