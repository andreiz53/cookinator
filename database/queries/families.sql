-- name: CreateFamily :one
INSERT INTO families (
    name,
    created_by_user_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetFamilyByID :one
SELECT * FROM families
WHERE id = $1;

-- name: GetFamilies :many
SELECT * FROM families;

-- name: GetFamilyByUserID :one
SELECT * FROM families 
WHERE created_by_user_id = $1;

-- name: UpdateFamily :one
UPDATE families SET
    updated_at = NOW(),
    name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteFamily :exec
DELETE FROM families
WHERE id = $1;