-- +goose Up
CREATE TABLE pastebin (
                          id SERIAL NOT NULL PRIMARY KEY,
                          text VARCHAR NOT NULL,
                          alias VARCHAR NOT NULL UNIQUE,
                          alias_for_del VARCHAR NOT NULL UNIQUE,
                          only_one BOOLEAN NOT NULL DEFAULT FALSE,
                          created_at DATE NOT NULL DEFAULT CURRENT_DATE,
                          deleted_at DATE
);
-- +goose Down
DROP TABLE pastebin;

