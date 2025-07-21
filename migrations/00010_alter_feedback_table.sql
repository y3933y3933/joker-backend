-- +goose Up
-- +goose StatementBegin
-- Add new review_status column
ALTER TABLE feedback 
ADD COLUMN review_status VARCHAR(20) NOT NULL DEFAULT 'new' 
CHECK (review_status IN ('new', 'reviewed'));

-- Migrate existing data: is_reviewed = true -> 'reviewed', false -> 'new'
UPDATE feedback 
SET review_status = CASE 
    WHEN is_reviewed = true THEN 'reviewed'
    ELSE 'new'
END;

-- Drop old column
ALTER TABLE feedback DROP COLUMN is_reviewed;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Add back is_reviewed column
ALTER TABLE feedback 
ADD COLUMN is_reviewed BOOLEAN NOT NULL DEFAULT FALSE;

-- Migrate data back: 'reviewed' -> true, others -> false
UPDATE feedback 
SET is_reviewed = CASE 
    WHEN review_status = 'reviewed' THEN true
    ELSE false
END;

-- Drop new column
ALTER TABLE feedback DROP COLUMN review_status;
-- +goose StatementEnd