-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  sender_id VARCHAR(92) NOT NULL,
  conversation_id INT NOT NULL,
  message VARCHAR(255) NOT NULL,
  FOREIGN KEY (sender_id) REFERENCES users (id),
  FOREIGN KEY (conversation_id) REFERENCES conversations (id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE messages;