-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA shortener;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS shortener;
-- +goose StatementEnd
