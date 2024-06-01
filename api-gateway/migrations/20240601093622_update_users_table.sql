-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE users
ADD COLUMN IF NOT EXISTS last_code_update TIMESTAMP;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE users
DROP COLUMN last_code_update TIMESTAMP;