-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE users
DROP COLUMN password,
ADD COLUMN code VARCHAR(10) NOT NULL;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users
DROP COLUMN code,
ADD COLUMN  password VARCHAR(32) NOT NULL;
-- +goose StatementEnd
