-- +migrate Up
CREATE TABLE boards (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
);

-- +migrate Down
DROP TABLE boards;
