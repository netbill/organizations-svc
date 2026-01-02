-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE outbox_event_status AS ENUM (
    'pending',
    'processing',
    'sent',
    'failed'
);

CREATE TABLE outbox_events (
    id       UUID  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    topic    TEXT  NOT NULL,
    key      TEXT  NOT NULL,
    type     TEXT  NOT NULL,
    version  INT   NOT NULL,
    producer TEXT  NOT NULL,
    payload  JSONB NOT NULL,

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
    id       UUID  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    topic    TEXT  NOT NULL,
    key      TEXT  NOT NULL,
    type     TEXT  NOT NULL,
    version  INT   NOT NULL,
    producer TEXT  NOT NULL,
    payload  JSONB NOT NULL,

    status        inbox_event_status NOT NULL DEFAULT 'pending', -- pending | processed | failed
    attempts      INT         NOT NULL DEFAULT 0,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    next_retry_at TIMESTAMPTZ,
    processed_at  TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS inbox_events CASCADE;

DROP TYPE IF EXISTS outbox_event_status;
DROP TYPE IF EXISTS inbox_event_status;