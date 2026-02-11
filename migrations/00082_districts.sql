-- +gooseUp

-- +goose StatementBegin
CREATE TABLE districts(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	name VARCHAR(255) NOT NULL,
	code VARCHAR(20) NOT NULL,
	country_id UUID NOT NULL REFERENCES countries(id),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_districts_table_version_code UNIQUE (table_code, version, code)
);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS districts;
-- +goose StatementEnd
