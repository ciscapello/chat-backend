-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE users ALTER COLUMN username DROP NOT NULL;
ALTER TABLE users
  ADD CONSTRAINT users_username_key UNIQUE (username);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
