-- +goose Up
-- +goose StatementBegin
CREATE TABLE shortener.urlmapping (
  slug        VARCHAR(8)    PRIMARY KEY,
  original    VARCHAR(2048) UNIQUE NOT NULL,
  user_id     UUID          NULL,
  created_at  TIMESTAMP     NOT NULL,
  expires_at  TIMESTAMP     NULL,
  deleted     boolean       NOT NULL
);
CREATE INDEX idx_user_id ON shortener.urlmapping (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS shortener.urlmapping; 
-- +goose StatementEnd
