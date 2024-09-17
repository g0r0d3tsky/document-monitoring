-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "users" (
                         "id" UUID PRIMARY KEY,
                         "userName" VARCHAR(255) UNIQUE NOT NULL,
                         "firstName" VARCHAR(255) NOT NULL,
                         "lastName" VARCHAR(255) NOT NULL,
                         "email" VARCHAR(255) UNIQUE NOT NULL,
                         "password" VARCHAR(255) NOT NULL,
                         "role" VARCHAR(50) NOT NULL DEFAULT 'user'
);

-- Создание таблицы files
CREATE TABLE "files" (
                         "id" UUID PRIMARY KEY,
                         "filename" VARCHAR(255) NOT NULL,
                         "file_type" VARCHAR(50) NOT NULL,
                         "file_size" BIGINT NOT NULL,
                         "file_path" VARCHAR(1024) NOT NULL,
                         "checksum" VARCHAR(255),
                         "storage_location" VARCHAR(100) NOT NULL DEFAULT 'local',
                         "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         "user_id" UUID NOT NULL
);

-- Добавление внешнего ключа
ALTER TABLE "files" ADD CONSTRAINT "files_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "files" DROP CONSTRAINT "files_user_id_fkey";

DROP TABLE IF EXISTS "files";
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
