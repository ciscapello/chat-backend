-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE role AS ENUM ('admin', 'regular');

CREATE TABLE users (
  id VARCHAR(92) NOT NULL PRIMARY KEY,
  username VARCHAR(32) NOT NULL UNIQUE,
  password VARCHAR(32) NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  role role NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;

DROP TYPE role;
-- +goose StatementEnd
