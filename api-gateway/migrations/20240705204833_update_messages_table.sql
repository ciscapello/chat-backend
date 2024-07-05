-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE messages
ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

UPDATE messages
SET created_at = NOW()
WHERE created_at IS NULL;

UPDATE messages
SET updated_at = NOW()
WHERE updated_at IS NULL;


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE messages
DROP COLUMN created_at;
ALTER TABLE messages
DROP COLUMN updated_at;
-- +goose StatementEnd
