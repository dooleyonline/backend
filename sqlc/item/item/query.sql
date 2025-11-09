-- name: GetAll :many
SELECT
  *
FROM
  item.item AS i
ORDER BY 
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'asc'  THEN i.posted_at END ASC,
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'desc' THEN i.posted_at END DESC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'asc'      THEN i.price END ASC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'desc'     THEN i.price END DESC, 

  posted_at DESC
LIMIT 
  CAST(@size AS integer) OFFSET CAST(@size AS integer)*(CAST(@page AS integer)-1);

-- name: Get :one
SELECT
  *
FROM
  item.item
WHERE
  id = $1;

-- name: Create :one
INSERT INTO
  item.item (name, description, images, price, condition, is_negotiable, category, subcategory, placeholder, seller)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: Update :one
UPDATE
  item.item
SET
  name = $2,
  description = $3,
  images = $4,
  price = $5,
  condition = $6,
  is_negotiable = $7,
  category = $8,
  subcategory = $9,
  placeholder = $10
WHERE
  id = $1
RETURNING *;

-- name: Sell :exec
UPDATE
  item.item
SET
  sold_at = $2
WHERE
  id = $1;

-- name: IncrementView :exec
UPDATE
  item.item
SET
  views = views + 1
WHERE
  id = $1;

-- name: IncrementLike :exec
UPDATE
  item.item
SET
  likes = likes + 1
WHERE
  id = $1; 

-- name: DecrementLike :exec
UPDATE
  item.item
SET
  likes = GREATEST(likes - 1, 0)
WHERE
  id = $1;

-- name: Delete :exec
DELETE FROM
  item.item
WHERE
  id = $1;

-- name: Search :many
SELECT
 i.*
FROM
  item.item AS i
WHERE
  i.fts @@ websearch_to_tsquery($1)
ORDER BY 
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'asc'  THEN i.posted_at END ASC,
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'desc' THEN i.posted_at END DESC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'asc'      THEN i.price END ASC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'desc'     THEN i.price END DESC, 

  ts_rank(i.fts, websearch_to_tsquery(@q)) DESC
LIMIT 
  CAST(@size AS integer) OFFSET CAST(@size AS integer)*(CAST(@page AS integer)-1);

-- name: GetByCategory :many
SELECT
  *
FROM
  item.item AS i 
WHERE
  category = $1
ORDER BY 
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'asc'  THEN i.posted_at END ASC,
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'desc' THEN i.posted_at END DESC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'asc'      THEN i.price END ASC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'desc'     THEN i.price END DESC, 

  posted_at DESC
LIMIT 
  CAST(@size AS integer) OFFSET CAST(@size AS integer)*(CAST(@page AS integer)-1);

-- name: SearchByCategory :many 
SELECT
  *
FROM
  item.item AS i
WHERE
  category = $1 AND fts @@ to_tsquery($2)
ORDER BY 
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'asc'  THEN i.posted_at END ASC,
  CASE WHEN @order_by = 'posted_at' AND @order_dir = 'desc' THEN i.posted_at END DESC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'asc'      THEN i.price END ASC,
  CASE WHEN @order_by = 'price' AND @order_dir = 'desc'     THEN i.price END DESC,
  
  ts_rank(i.fts, websearch_to_tsquery(@q)) DESC
LIMIT 
  CAST(@size AS integer) OFFSET CAST(@size AS integer)*(CAST(@page AS integer)-1);

-- name: GetBySeller :many
SELECT
  *
FROM
  item.item
WHERE
  seller = @seller_ID::uuid
ORDER BY
  posted_at DESC
LIMIT 
  CAST(@size AS integer) OFFSET CAST(@size AS integer)*(CAST(@page AS integer)-1);

-- name: GetBatch :many 
SELECT
  *
FROM
  item.item
WHERE
  id = ANY(@item_IDs::bigint[])
ORDER BY
  posted_at DESC;


