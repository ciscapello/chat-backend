-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users
ADD COLUMN email VARCHAR(40)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
