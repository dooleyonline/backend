-- name: GetMany :many
SELECT
  *
FROM
  "user".verify;


-- name: Get :one
SELECT
  *
FROM
  "user".verify
WHERE
  id = $1; 


-- name: Create :one
INSERT INTO
  "user".verify (user_id)
VALUES
  ($1)
RETURNING *;

