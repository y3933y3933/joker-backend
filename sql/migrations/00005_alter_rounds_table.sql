-- +goose Up
-- +goose StatementBegin
ALTER TABLE rounds
ALTER COLUMN question_id DROP NOT NULL;

ALTER TABLE rounds
ADD COLUMN answer TEXT DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rounds
ALTER COLUMN question_id SET NOT NULL;

ALTER TABLE rounds
DROP COLUMN IF EXISTS answer;
-- +goose StatementEnd