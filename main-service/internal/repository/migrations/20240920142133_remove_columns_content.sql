-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "files"
    DROP COLUMN IF EXISTS "storage_location",
    DROP COLUMN IF EXISTS "file_type";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "files"
    ADD COLUMN IF NOT EXISTS "storage_location" VARCHAR(255),
    ADD COLUMN IF NOT EXISTS "file_type" VARCHAR(50);
-- +goose StatementEnd
