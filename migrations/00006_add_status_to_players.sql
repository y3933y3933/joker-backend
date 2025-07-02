-- +goose Up
-- +goose StatementBegin
ALTER TABLE players ADD COLUMN status TEXT NOT NULL DEFAULT 'online';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players DROP COLUMN status;
-- +goose StatementEnd
