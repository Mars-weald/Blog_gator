-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR,
    url VARCHAR NOT NULL UNIQUE,
    description VARCHAR,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;