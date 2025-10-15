-- name: GetAll :many 
SELECT
  *
FROM
  "user";

-- name: Create :one
INSERT INTO
  "user" (email, password)
VALUES
  ($1, $2)
RETURNING *;

-- name: Get :one
SELECT
  *
FROM
  "user"
WHERE
  email = $1;