-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shortener.urlmapping (
  slug        VARCHAR(8)    PRIMARY KEY,
  original    VARCHAR(2048) NOT NULL,
  created_at  TIMESTAMP     NOT NULL,
  expires_at  TIMESTAMP     NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shortener.urlmapping; 
-- +goose StatementEnd