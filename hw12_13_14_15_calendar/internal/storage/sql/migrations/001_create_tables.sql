-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.events (
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    timestamp_event timestamp without time zone NOT NULL,
    duration INTEGER NOT NULL,
    description VARCHAR,
    customer_id uuid NOT NULL,
    notify_duration INTEGER
);
