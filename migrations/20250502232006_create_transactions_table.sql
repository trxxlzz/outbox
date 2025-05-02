-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
                              id VARCHAR(255) PRIMARY KEY,
                              user_id VARCHAR(255) NOT NULL,
                              amount NUMERIC(20, 2) NOT NULL,
                              currency VARCHAR(10) NOT NULL,
                              status VARCHAR(50) NOT NULL,
                              timestamp TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
