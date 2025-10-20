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
  item (name, description, images, price, condition, is_negotiable, category, subcategory, placeholder, seller)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
  category = $8,
  subcategory = $9,
  views = $10,
  placeholder = $11
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

-- name: IncrementLike :exec
UPDATE
  item
SET
  likes = likes + 1
WHERE
  id = $1; 

-- name: DecrementLike :exec
UPDATE
  item
SET
  likes = GREATEST(likes - 1, 0)
WHERE
  id = $1;

-- name: Delete :exec
DELETE FROM
  item
WHERE
  id = $1;

-- name: Search :many
SELECT
 *
FROM
  item
WHERE
  fts @@ websearch_to_tsquery($1);

-- name: GetByCategory :many
SELECT
  *
FROM
  item
WHERE
  category = $1;

-- name: SearchByCategory :many 
SELECT
  *
FROM
  item
WHERE
  category = $1 AND fts @@ to_tsquery($2); 

-- name: GetBySeller :many
SELECT
  *
FROM
  "item"
WHERE
  seller = @seller_ID::uuid
ORDER BY
  posted_at DESC;

-- name: GetByIDs :many 
SELECT
  *
FROM
  "item"
WHERE
  id = ANY(@item_IDs::bigint[]);