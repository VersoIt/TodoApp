-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_lists
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS users_lists
(
    id      SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    list_id INT REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_items
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    done        BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS lists_items
(
    id      SERIAL PRIMARY KEY,
    item_id INT REFERENCES todo_items (id) ON DELETE CASCADE NOT NULL,
    list_id INT REFERENCES todo_lists (id) ON DELETE CASCADE  NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lists_items;
DROP TABLE IF EXISTS users_lists;
DROP TABLE IF EXISTS todo_lists;
DROP TABLE IF EXiSTS users;
DROP TABLE IF EXISTS todo_items;
-- +goose StatementEnd
