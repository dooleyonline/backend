-- name: GetMany :many
SELECT
  *
FROM
  "user"."user";

-- name: Get :one
SELECT
  *
FROM
  "user"."user"
WHERE
  email = $1;

-- name: Create :one
INSERT INTO
  "user"."user" (email, password, first_name, last_name)
VALUES
  ($1, $2, $3, $4)
RETURNING *;

-- name: GetLikedItems :one 
SELECT
  liked_items
FROM 
  "user"."user"
WHERE
  id = $1;

-- name: DeleteLikedItem :exec
UPDATE
  "user"."user"
SET
  liked_items = array_remove(liked_items, @item_ID::bigint)
WHERE
  id = @id
  AND (liked_items @> ARRAY[@item_ID::bigint]);

-- name: AddLikedItem :exec
UPDATE
  "user"."user"
SET
  liked_items = array_append(liked_items, @item_ID::bigint)
WHERE
  id = @id
  AND NOT liked_items @> ARRAY[@item_ID::bigint];

-- name: GetByID :one
SELECT
  *
FROM
  "user"."user"
WHERE
  id = $1;

-- name: Update :exec
UPDATE
  "user"."user"
SET
  email = $2,
  first_name = $3,
  last_name = $4,
  avatar = $5
WHERE
  id = $1;


-- name: Verify :exec
UPDATE
  "user"."user"
SET
  verified = TRUE
WHERE
  id = $1;
