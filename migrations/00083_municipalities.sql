-- +gooseUp

-- +goose StatementBegin
CREATE TABLE municipalities(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	name VARCHAR(255) NOT NULL,
	code VARCHAR(20) NOT NULL,
	district_id UUID NOT NULL REFERENCES districts(id),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_municipalities_table_version_code UNIQUE (table_code, version, code)
);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS municipalities;
-- +goose StatementEnd
