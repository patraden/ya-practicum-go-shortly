-- +goose Up
-- +goose StatementBegin
CREATE TABLE shortener.urlmapping (
  slug        VARCHAR(8)    PRIMARY KEY,
  original    VARCHAR(2048) UNIQUE NOT NULL,
  created_at  TIMESTAMP     NOT NULL,
  expires_at  TIMESTAMP     NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shortener.urlmapping; 
-- +goose StatementEnd
