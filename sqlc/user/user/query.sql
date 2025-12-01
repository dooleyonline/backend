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

-- name: GetAllLikedViewed :many
SELECT
    u.user_id,
    ARRAY_AGG(DISTINCT l.item_id)  FILTER (WHERE l.item_id IS NOT NULL)  AS liked_items,
    ARRAY_AGG(DISTINCT v.item_id) FILTER (WHERE v.item_id IS NOT NULL) AS viewed_items
FROM (
    SELECT user_id FROM "user"."liked"
    UNION
    SELECT user_id FROM "user"."viewed"
) AS u
LEFT JOIN "user"."liked"  AS l ON l.user_id = u.user_id
LEFT JOIN "user"."viewed" AS v ON v.user_id = u.user_id
GROUP BY u.user_id
ORDER BY u.user_id;
