-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id    UUID PRIMARY KEY,
    username   TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

