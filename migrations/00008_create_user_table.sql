-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
     id BIGSERIAL PRIMARY KEY,
     username TEXT NOT NULL UNIQUE,
     password_hash bytea NOT NULL,
     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
