-- name: CreateIngredient :one
INSERT INTO ingredients (
    name, 
    density
) VALUES ( $1, $2)
RETURNING *;

-- name: GetIngredientByID :one
SELECT * FROM ingredients
WHERE id = $1;

-- name: GetIngredientByName :one
SELECT * FROM ingredients
WHERE name = $1;

-- name: GetIngredients :many
SELECT * FROM ingredients;

-- name: UpdateIngredient :one
UPDATE ingredients SET
    name = $2,
    density = $3
WHERE id = $1
RETURNING *;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1;