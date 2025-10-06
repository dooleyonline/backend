-- name: GetAll :many
SELECT
 id, name, description, images, price, condition, is_negotiable, posted_at, sold_at, views, category, sub_category
FROM
 item;

-- name: Get :one
SELECT
  id, name, description, images, price, condition, is_negotiable, posted_at, sold_at, views, category, sub_category
FROM
  item
WHERE
  id = $1;

-- name: Create :one
INSERT INTO
  item (name, description, images, price, condition, is_negotiable, category, sub_category)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, name, description, images, price, condition, is_negotiable, posted_at, sold_at, views, category, sub_category;

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
RETURNING id, name, description, images, price, condition, is_negotiable, posted_at, sold_at, views, category, sub_category;

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

-- name: Search :many
SELECT
  id, name, description, images, price, condition, is_negotiable, posted_at, sold_at, views, category, sub_category
FROM
  item
WHERE
  fts @@ to_tsquery($1);
