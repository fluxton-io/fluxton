-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.form_responses (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uuid UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_form_responses_form FOREIGN KEY (form_uuid) REFERENCES fluxton.forms(uuid) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fluxton.form_responses CASCADE;
-- +goose StatementEnd
