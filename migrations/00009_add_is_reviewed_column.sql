-- +goose Up
-- +goose StatementBegin
ALTER TABLE feedback 
ADD COLUMN is_reviewed BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feedback DROP COLUMN is_reviewed;
-- +goose StatementEnd
