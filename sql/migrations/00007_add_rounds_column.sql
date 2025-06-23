-- +goose Up
-- +goose StatementBegin
ALTER TABLE rounds ADD COLUMN deck TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rounds DROP COLUMN deck;
-- +goose StatementEnd
