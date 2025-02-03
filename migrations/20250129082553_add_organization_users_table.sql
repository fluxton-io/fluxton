-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_users (
    id SERIAL PRIMARY KEY,
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE organization_users ADD CONSTRAINT fk_organization_users_organizations FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE;
ALTER TABLE organization_users ADD CONSTRAINT fk_organization_users_users FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_users;
-- +goose StatementEnd
