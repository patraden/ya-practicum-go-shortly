-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS shortener;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS shortener;
-- +goose StatementEnd
