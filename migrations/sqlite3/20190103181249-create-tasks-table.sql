-- +migrate Up
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    estimate INT DEFAULT 0,
    spent INT DEFAULT 0
);

-- +migrate Down
DROP TABLE tasks;
