-- +goose Up
-- +goose StatementBegin
ALTER TABLE test ADD COLUMN user_id uuid;
ALTER TABLE test ADD COLUMN balance decimal(10, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test DROP COLUMN user_id;
ALTER TABLE test DROP COLUMN balance;
-- +goose StatementEnd
