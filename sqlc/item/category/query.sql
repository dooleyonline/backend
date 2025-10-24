-- name: GetAll :many 
SELECT
  *
FROM
  item.category;

-- name: Create :one
INSERT INTO
  item.category (name, subcategory, icon)
VALUES
  ($1, $2, $3)
RETURNING *;

-- name: Get :one
SELECT
  *
FROM
  item.category
WHERE
  name = $1;
