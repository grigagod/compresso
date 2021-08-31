CREATE SCHEMA IF NOT EXISTS svc AUTHORIZATION postgres;

-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS svc.users (
    user_id    UUID PRIMARY KEY,
    username   TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'image_format') THEN
        CREATE TYPE image_format AS ENUM('jpg', 'png');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing') THEN
        CREATE TYPE processing_state as ENUM('queued', 'processing', 'done', 'failed');
    END IF;
END$$;

--CREATE TYPE IF NOT EXISTS image_format as ENUM('jpg', 'png');

CREATE TABLE IF NOT EXISTS svc.image_tickets (
    ticket_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id     UUID NOT NULL REFERENCES svc.users (user_id),
    format        image_format NOT NULL,
    state         processing_state NOT NULL,
    compression   SMALLINT NOT NULL,
    origin_url    TEXT NOT NULL,
    processed_url TEXT,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
