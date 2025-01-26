-- name: CreateRecipe :one
INSERT INTO recipes (
    name,
    cooking_process,
    family_id,
    items
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetRecipes :many
SELECT * FROM recipes;

-- name: GetRecipesByFamilyID :many
SELECT * FROM recipes
WHERE family_id = $1;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1;

-- name: UpdateRecipe :one
UPDATE recipes SET
    name = $2,
    cooking_process = $3,
    items = $4
WHERE id = $1
RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;