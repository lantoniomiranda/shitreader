-- +gooseUp

-- +goose StatementBegin
CREATE TABLE process_suspension_reasons(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	code VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_process_suspension_reasons_table_version_code UNIQUE (table_code, version, code)
);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS process_suspension_reasons;
-- +goose StatementEnd
