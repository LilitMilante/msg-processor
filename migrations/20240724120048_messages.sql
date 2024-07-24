-- +goose Up
-- +goose StatementBegin
CREATE TYPE state AS ENUM ('UNPROCESSED', 'PROCESSED');

CREATE TABLE messages (
    id UUID NOT NULL,
    text TEXT NOT NULL,
    state state NOT NULL,
    created_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
