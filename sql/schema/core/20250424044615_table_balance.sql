-- +goose Up
-- +goose StatementBegin
CREATE TABLE balance (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    balance INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balance;
-- +goose StatementEnd
