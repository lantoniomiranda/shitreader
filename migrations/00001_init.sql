-- +gooseUp
-- +goose StatementBegin

-- 1. Table Versions
CREATE TABLE table_versions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_code VARCHAR(50) NOT NULL,
	version VARCHAR(20) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_table_version UNIQUE (table_code, version)
);

-- 2. Catalogs
CREATE TABLE catalogs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	slug VARCHAR(100) NOT NULL UNIQUE,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL
);

-- 3. Catalog Values
CREATE TABLE catalog_values (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	catalog_id UUID NOT NULL REFERENCES catalogs(id) ON DELETE CASCADE,
	table_version_id UUID NOT NULL REFERENCES table_versions(id) ON DELETE CASCADE,
	code VARCHAR(50) NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_catalog_value UNIQUE (catalog_id, table_version_id, code)
);

-- 4. Geography Tables
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_version_id UUID NOT NULL REFERENCES table_versions(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT unique_countries_table_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_version_id UUID NOT NULL REFERENCES table_versions(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    country_id UUID REFERENCES countries(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT unique_districts_table_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE municipalities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_version_id UUID NOT NULL REFERENCES table_versions(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    district_id UUID REFERENCES districts(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT unique_municipalities_table_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE parishes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_version_id UUID NOT NULL REFERENCES table_versions(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    municipality_id UUID REFERENCES municipalities(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT unique_parishes_table_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE ine_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_version_id UUID NOT NULL REFERENCES table_versions(id),
    zone_code VARCHAR(20) NOT NULL,
    zone_name VARCHAR(255) NOT NULL,
    zone_name_formatted VARCHAR(255) NOT NULL,
    ine_municipality_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT unique_ine_zones_table_version_code UNIQUE (table_version_id, zone_code)
);

-- 5. Structural Tables
CREATE TABLE steps (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_version_id UUID NOT NULL REFERENCES table_versions(id),
	code VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_steps_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE records (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_version_id UUID NOT NULL REFERENCES table_versions(id),
	code VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,
	record_type_id UUID REFERENCES catalog_values(id),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_records_version_code UNIQUE (table_version_id, code)
);

CREATE TABLE fields (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	table_version_id UUID NOT NULL REFERENCES table_versions(id),
	code VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,
	record_id UUID REFERENCES records(id),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_fields_version_code UNIQUE (table_version_id, code)
);

-- 6. Association Tables
CREATE TABLE step_header_types(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	step_id UUID NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
	header_type_id UUID NOT NULL REFERENCES catalog_values(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_step_header_type UNIQUE (step_id, header_type_id)
);

CREATE INDEX idx_step_header_types_step_id ON step_header_types(step_id);
CREATE INDEX idx_step_header_types_header_type_id ON step_header_types(header_type_id);

CREATE TABLE step_records(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	step_id UUID NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
	record_id UUID NOT NULL REFERENCES records(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_step_record UNIQUE (step_id, record_id)
);

CREATE INDEX idx_step_records_step_id ON step_records(step_id);
CREATE INDEX idx_step_records_record_id ON step_records(record_id);

-- 7. Process Tables
CREATE TABLE processes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(20) NOT NULL UNIQUE,
	description TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE process_steps (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	process_id UUID NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
	step_id UUID NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
	step_order INT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT unique_process_step UNIQUE (process_id, step_id)
);

CREATE INDEX idx_process_steps_process_id ON process_steps(process_id);
CREATE INDEX idx_process_steps_step_id ON process_steps(step_id);

-- +goose StatementEnd

-- +gooseDown
-- +goose StatementBegin
DROP TABLE IF EXISTS process_steps CASCADE;
DROP TABLE IF EXISTS processes CASCADE;
DROP TABLE IF EXISTS step_records CASCADE;
DROP TABLE IF EXISTS step_header_types CASCADE;
DROP TABLE IF EXISTS fields CASCADE;
DROP TABLE IF EXISTS records CASCADE;
DROP TABLE IF EXISTS steps CASCADE;
DROP TABLE IF EXISTS ine_zones CASCADE;
DROP TABLE IF EXISTS parishes CASCADE;
DROP TABLE IF EXISTS municipalities CASCADE;
DROP TABLE IF EXISTS districts CASCADE;
DROP TABLE IF EXISTS countries CASCADE;
DROP TABLE IF EXISTS catalog_values CASCADE;
DROP TABLE IF EXISTS catalogs CASCADE;
DROP TABLE IF EXISTS table_versions CASCADE;
-- +goose StatementEnd
