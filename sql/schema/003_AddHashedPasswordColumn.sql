-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset'; --Default is totally useless and only making db worse but stays because its a guiding project

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;
