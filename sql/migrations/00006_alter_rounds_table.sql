-- +goose Up
-- +goose StatementBegin
ALTER TABLE rounds
RENAME COLUMN current_player_id TO question_player_id;

ALTER TABLE rounds
ADD COLUMN answer_player_id BIGINT;

ALTER TABLE rounds
ADD CONSTRAINT fk_answer_player
FOREIGN KEY (answer_player_id) REFERENCES players(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rounds DROP CONSTRAINT fk_answer_player;

ALTER TABLE rounds DROP COLUMN answer_player_id;

ALTER TABLE rounds
RENAME COLUMN question_player_id TO current_player_id;
-- +goose StatementEnd
