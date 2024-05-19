-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users
ADD CONSTRAINT unique_email UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
