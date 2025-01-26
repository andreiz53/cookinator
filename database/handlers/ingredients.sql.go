// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: ingredients.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createIngredient = `-- name: CreateIngredient :one
INSERT INTO ingredients (
    name, 
    density
) VALUES ( $1, $2)
RETURNING id, name, density
`

type CreateIngredientParams struct {
	Name    string
	Density pgtype.Numeric
}

func (q *Queries) CreateIngredient(ctx context.Context, arg CreateIngredientParams) (Ingredient, error) {
	row := q.db.QueryRow(ctx, createIngredient, arg.Name, arg.Density)
	var i Ingredient
	err := row.Scan(&i.ID, &i.Name, &i.Density)
	return i, err
}

const deleteIngredient = `-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1
`

func (q *Queries) DeleteIngredient(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteIngredient, id)
	return err
}

const getIngredientByID = `-- name: GetIngredientByID :one
SELECT id, name, density FROM ingredients
WHERE id = $1
`

func (q *Queries) GetIngredientByID(ctx context.Context, id int32) (Ingredient, error) {
	row := q.db.QueryRow(ctx, getIngredientByID, id)
	var i Ingredient
	err := row.Scan(&i.ID, &i.Name, &i.Density)
	return i, err
}

const getIngredientByName = `-- name: GetIngredientByName :one
SELECT id, name, density FROM ingredients
WHERE name = $1
`

func (q *Queries) GetIngredientByName(ctx context.Context, name string) (Ingredient, error) {
	row := q.db.QueryRow(ctx, getIngredientByName, name)
	var i Ingredient
	err := row.Scan(&i.ID, &i.Name, &i.Density)
	return i, err
}

const getIngredients = `-- name: GetIngredients :many
SELECT id, name, density FROM ingredients
`

func (q *Queries) GetIngredients(ctx context.Context) ([]Ingredient, error) {
	rows, err := q.db.Query(ctx, getIngredients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ingredient
	for rows.Next() {
		var i Ingredient
		if err := rows.Scan(&i.ID, &i.Name, &i.Density); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateIngredient = `-- name: UpdateIngredient :one
UPDATE ingredients SET
    name = $2,
    density = $3
WHERE id = $1
RETURNING id, name, density
`

type UpdateIngredientParams struct {
	ID      int32
	Name    string
	Density pgtype.Numeric
}

func (q *Queries) UpdateIngredient(ctx context.Context, arg UpdateIngredientParams) (Ingredient, error) {
	row := q.db.QueryRow(ctx, updateIngredient, arg.ID, arg.Name, arg.Density)
	var i Ingredient
	err := row.Scan(&i.ID, &i.Name, &i.Density)
	return i, err
}
