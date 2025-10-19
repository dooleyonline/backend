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
  "user" (email, password, first_name, last_name)
VALUES
  ($1, $2, $3, $4)
RETURNING *;

-- name: GetLikedItems :one 
SELECT
  liked_items
FROM 
  "user"
WHERE
  email = $1;

-- name: DeleteLikedItem :one
UPDATE
  "user"
SET
  liked_items = array_remove(liked_items, @item_ID::bigint)
WHERE
  email = @email
  AND (liked_items @> ARRAY[@item_ID::bigint])
RETURNING liked_items; 

-- name: AddLikedItem :one 
UPDATE
  "user"
SET
  liked_items = array_append(liked_items, @item_ID::bigint)
WHERE
  email = @email
  AND NOT liked_items @> ARRAY[@item_ID::bigint]
RETURNING liked_items;