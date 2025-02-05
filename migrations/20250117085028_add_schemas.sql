-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS fluxton;
CREATE SCHEMA IF NOT EXISTS authentication;
CREATE SCHEMA IF NOT EXISTS storage;

SET search_path TO public, fluxton, authentication, storage;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS fluxton CASCADE;
DROP SCHEMA IF EXISTS authentication CASCADE;
DROP SCHEMA IF EXISTS storage CASCADE;
-- Avoid dropping the public schema as it's required for system tables
-- +goose StatementEnd
