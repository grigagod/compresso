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
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'video_format') THEN
        CREATE TYPE image_format AS ENUM('mkv', 'webm');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing') THEN
        CREATE TYPE processing_state as ENUM('queued', 'processing', 'done', 'failed');
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS svc.user_videos (
    video_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id     UUID NOT NULL REFERENCES svc.users (user_id),
    format        video_format NOT NULL,
    url           TEXT NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);


CREATE TABLE IF NOT EXISTS svc.video_tickets (
    ticket_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id      UUID NOT NULL REFERENCES svc.user_videos (video_id),
    author_id     UUID NOT NULL REFERENCES svc.users (user_id),
    state         processing_state NOT NULL,
    format        video_format NOT NULL,
    crf           SMALLINT NOT NULL,
    url           TEXT,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
