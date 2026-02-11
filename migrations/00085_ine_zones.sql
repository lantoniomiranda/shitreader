-- +gooseUp

-- +goose StatementBegin
CREATE TABLE ine_zones(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	zone_code VARCHAR(20) NOT NULL,
	zone_name VARCHAR(255) NOT NULL,
	zone_name_formatted VARCHAR(255) NOT NULL,
	ine_municipality_code VARCHAR(20) NOT NULL, 
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_ine_zones_table_version_code UNIQUE (table_code, version, zone_code)
);
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS ine_zones;
-- +goose StatementEnd
