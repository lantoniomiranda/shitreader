-- +gooseUp

-- +goose StatementBegin
CREATE TABLE steps(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	code VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,
	header_type_id UUID REFERENCES header_types(id),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_steps_table_version_code UNIQUE (table_code, version, code)
);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS steps;
-- +goose StatementEnd
