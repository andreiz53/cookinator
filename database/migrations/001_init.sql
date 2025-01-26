-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    first_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL
);

ALTER TABLE users ADD COLUMN family_id UUID REFERENCES families(id) ON DELETE SET NULL;

CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    density NUMERIC NOT NULL
);

CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    cooking_process TEXT NOT NULL,
    family_id UUID REFERENCES families(id) ON DELETE CASCADE,
    items JSONB NOT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_family_id ON users(family_id);

CREATE UNIQUE INDEX idx_families_name ON families(name);
CREATE INDEX idx_families_created_by ON families(created_by_user_id);

CREATE UNIQUE INDEX idx_ingredients_name ON ingredients(name);

CREATE INDEX idx_recipes_family_id ON recipes(family_id);
CREATE INDEX idx_recipes_name ON recipes(name);
CREATE INDEX idx_recipes_items ON recipes USING GIN (items);


-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_family_id_fkey;
ALTER TABLE families DROP CONSTRAINT IF EXISTS families_created_by_user_id_fkey;
ALTER TABLE recipes DROP CONSTRAINT IF EXISTS recipes_family_id_fkey;

DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS families;