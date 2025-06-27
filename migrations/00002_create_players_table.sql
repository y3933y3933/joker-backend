-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS players (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGSERIAL NOT NULL, 
    nickname VARCHAR(100) NOT NULL,
    is_host BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY(game_id) REFERENCES games(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
