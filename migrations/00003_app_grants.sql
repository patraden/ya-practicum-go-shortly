-- +goose Up
-- +goose StatementBegin
GRANT ALL PRIVILEGES ON SCHEMA shortener TO postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE ALL PRIVILEGES ON SCHEMA shortener FROM postgres;
-- +goose StatementEnd
