-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE conversations (
  id SERIAL PRIMARY KEY,
  creator_id VARCHAR(92) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY (creator_id) REFERENCES users (id)
);

CREATE TABLE participants (
  id SERIAL PRIMARY KEY,
  conversation_id INT NOT NULL,
  user_id VARCHAR(92) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY (conversation_id) REFERENCES conversations (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE conversations;
DROP TABLE participants;