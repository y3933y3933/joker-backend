-- +goose Up
ALTER TABLE rounds
ALTER COLUMN question_id DROP NOT NULL;

-- +goose Down
ALTER TABLE rounds
ALTER COLUMN question_id SET NOT NULL;
