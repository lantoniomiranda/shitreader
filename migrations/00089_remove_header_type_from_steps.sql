-- +gooseUp

-- +goose StatementBegin
INSERT INTO step_header_types (step_id, header_type_id)
SELECT id, header_type_id
FROM steps
WHERE header_type_id IS NOT NULL
ON CONFLICT (step_id, header_type_id) DO NOTHING;

ALTER TABLE steps
	DROP COLUMN header_type_id;
-- +goose StatementEnd

-- +gooseDown

-- +goose StatementBegin
ALTER TABLE steps
	ADD COLUMN header_type_id UUID REFERENCES header_types(id);

WITH ranked AS (
	SELECT
		step_id,
		header_type_id,
		ROW_NUMBER() OVER (PARTITION BY step_id ORDER BY created_at, header_type_id) AS rn
	FROM step_header_types
	WHERE deleted_at IS NULL
)
UPDATE steps s
SET header_type_id = r.header_type_id
FROM ranked r
WHERE s.id = r.step_id
	AND r.rn = 1;
-- +goose StatementEnd
