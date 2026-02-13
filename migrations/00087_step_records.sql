-- +gooseUp

-- +goose StatementBegin
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
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
DROP TABLE IF EXISTS step_records;
-- +goose StatementEnd
