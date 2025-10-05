-- name: GetAll :many 
SELECT
  *
FROM
  category;

-- name: Create :one
INSERT INTO
  category (name, subcategory, icon)
VALUES
  ($1, $2, $3)
RETURNING *;

