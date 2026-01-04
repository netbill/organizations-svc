-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE outbox_event_status AS ENUM (
    'pending',
    'processing',
    'sent',
    'failed'
);

CREATE TABLE outbox_events (
    id       UUID   PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    seq      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    topic    TEXT   NOT NULL,
    key      TEXT   NOT NULL,
    type     TEXT   NOT NULL,
    version  INT    NOT NULL,
    producer TEXT   NOT NULL,
    payload  JSONB  NOT NULL,

    status        outbox_event_status NOT NULL DEFAULT 'pending', -- pending | sent | failed
    attempts      INT         NOT NULL DEFAULT 0,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    next_retry_at TIMESTAMPTZ,
    sent_at       TIMESTAMPTZ
);

CREATE TYPE inbox_event_status AS ENUM (
    'pending',
    'processing',
    'processed',
    'failed'
);

CREATE TABLE inbox_events (
    id       UUID   PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    seq      BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    topic    TEXT   NOT NULL,
    key      TEXT   NOT NULL,
    type     TEXT   NOT NULL,
    version  INT    NOT NULL,
    producer TEXT   NOT NULL,
    payload  JSONB  NOT NULL,

    status        inbox_event_status NOT NULL DEFAULT 'pending', -- pending | processed | failed
    attempts      INT         NOT NULL DEFAULT 0,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    next_retry_at TIMESTAMPTZ,
    processed_at  TIMESTAMPTZ
);

CREATE INDEX idx_outbox_events_pending
    ON outbox_events (status, next_retry_at, seq);

CREATE INDEX idx_outbox_events_seq
    ON outbox_events (seq);

-- +migrate Down
DROP INDEX IF EXISTS idx_outbox_events_seq;
DROP INDEX IF EXISTS idx_outbox_events_pending;

DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS inbox_events CASCADE;

DROP TYPE IF EXISTS outbox_event_status;
DROP TYPE IF EXISTS inbox_event_status;