-- +gooseUp

-- +goose StatementBegin
CREATE TABLE step_header_types(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	step_id UUID NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
	header_type_id UUID NOT NULL REFERENCES header_types(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_step_header_type UNIQUE (step_id, header_type_id)
);

CREATE INDEX idx_step_header_types_step_id ON step_header_types(step_id);
CREATE INDEX idx_step_header_types_header_type_id ON step_header_types(header_type_id);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS step_header_types;
-- +goose StatementEnd
