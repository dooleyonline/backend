-- name: GetAll :many
SELECT
 *
FROM
 item;

-- name: Get :one
SELECT
  *
FROM
  item
WHERE
  id = $1;

-- name: Create :one
INSERT INTO
  item (name, description, images, price, condition, is_negotiable)
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: Update :one
UPDATE
  item
SET
  name = $2,
  description = $3,
  images = $4,
  price = $5,
  condition = $6,
  is_negotiable = $7,
  views = $8
WHERE
  id = $1
RETURNING *;

-- name: Sell :exec
UPDATE
  item
SET
  sold_at = $2
WHERE
  id = $1;

-- name: IncrementView :exec
UPDATE
  item
SET
  views = views + 1
WHERE
  id = $1;

-- name: Delete :exec
DELETE FROM
  item
WHERE
  id = $1;

