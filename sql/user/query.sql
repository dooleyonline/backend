-- name: GetMany :many
SELECT
  *
FROM
  "user";

-- name: Get :one
SELECT
  *
FROM
  "user"
WHERE
  email = $1;

-- name: Create :one
INSERT INTO
  "user" (email, password)
VALUES
  ($1, $2)
RETURNING *;
