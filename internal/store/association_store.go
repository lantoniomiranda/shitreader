package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type PostgresAssociationStore struct {
	db *sql.DB
}

func NewPostgresAssociationStore(db *sql.DB) *PostgresAssociationStore {
	return &PostgresAssociationStore{
		db: db,
	}
}

type AssociationStore interface {
	AssociateRecordsFields(ctx context.Context) error
	AssociateRecordsRecordTypes(ctx context.Context, filePath string, sheetName string) error
	AssociateStepsHeaderTypesAndRecords(ctx context.Context, filePath string, sheetName string) error
}

// AssociateRecordsFields links each field to its parent record using a single
// UPDATE ... FROM statement instead of per-row SELECT + UPDATE.
func (s *PostgresAssociationStore) AssociateRecordsFields(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE fields f
		SET record_id = r.id
		FROM records r
		WHERE LEFT(f.code, 5) = LEFT(r.code, 5)
		  AND f.deleted_at IS NULL
		  AND r.deleted_at IS NULL
		  AND f.record_id IS NULL
	`

	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to associate fields with records: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit field-record association: %w", err)
	}

	_ = result
	return nil
}

// AssociateRecordsRecordTypes updates records with their record_type_id
// based on data from an Excel file, wrapped in a single transaction.

func (s *PostgresAssociationStore) AssociateRecordsRecordTypes(ctx context.Context, filePath string, sheetName string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	// Load all record types from database into a map.
	recordTypesMap := make(map[string]string)
	recordTypesQuery := `
		SELECT cv.id, cv.code 
		FROM catalog_values cv
		JOIN catalogs c ON cv.catalog_id = c.id
		WHERE c.slug = 'record_types' AND cv.deleted_at IS NULL
	`

	recordTypeRows, err := s.db.QueryContext(ctx, recordTypesQuery)
	if err != nil {
		return fmt.Errorf("failed to load record types: %w", err)
	}
	defer recordTypeRows.Close()

	for recordTypeRows.Next() {
		var id, code string
		if err := recordTypeRows.Scan(&id, &code); err != nil {
			return fmt.Errorf("failed to scan record type: %w", err)
		}
		recordTypesMap[code] = id
	}

	// Begin a single transaction for all updates.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	associatedCount := 0
	skippedCount := 0

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}

		recordCode := row[1]
		recordTypeCode := row[2]

		if recordCode == "" || recordTypeCode == "" {
			continue
		}

		recordTypeId, exists := recordTypesMap[recordTypeCode]
		if !exists {
			skippedCount++
			continue
		}

		updateQuery := `
			UPDATE records 
			SET record_type_id = $1 
			WHERE code = $2 
			AND deleted_at IS NULL
		`

		result, err := tx.ExecContext(ctx, updateQuery, recordTypeId, recordCode)
		if err != nil {
			return fmt.Errorf("failed to update record %s with record_type_id %s: %w", recordCode, recordTypeId, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			skippedCount++
			continue
		}

		associatedCount++
		// No per-record logging to keep output minimal.
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit record-type association: %w", err)
	}

	return nil
}

// AssociateStepsHeaderTypesAndRecords creates step-header_type and step-record
// associations from an Excel file, wrapped in a single transaction.

func (s *PostgresAssociationStore) AssociateStepsHeaderTypesAndRecords(ctx context.Context, filePath string, sheetName string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	// Load all header types from database into a map.
	headerTypesMap := make(map[string]string)
	headerTypesQuery := `
		SELECT cv.id, cv.code 
		FROM catalog_values cv
		JOIN catalogs c ON cv.catalog_id = c.id
		WHERE c.slug = 'header_types' AND cv.deleted_at IS NULL
	`

	headerTypeRows, err := s.db.QueryContext(ctx, headerTypesQuery)
	if err != nil {
		return fmt.Errorf("failed to load header types: %w", err)
	}
	defer headerTypeRows.Close()

	for headerTypeRows.Next() {
		var id, code string
		if err := headerTypeRows.Scan(&id, &code); err != nil {
			return fmt.Errorf("failed to scan header type: %w", err)
		}
		headerTypesMap[code] = id
	}

	// Load all records from database into a map.
	recordsMap := make(map[string]string)
	recordsQuery := `SELECT id, code FROM records WHERE deleted_at IS NULL`

	recordRows, err := s.db.QueryContext(ctx, recordsQuery)
	if err != nil {
		return fmt.Errorf("failed to load records: %w", err)
	}
	defer recordRows.Close()

	for recordRows.Next() {
		var id, code string
		if err := recordRows.Scan(&id, &code); err != nil {
			return fmt.Errorf("failed to scan record: %w", err)
		}
		recordsMap[code] = id
	}

	// Load all steps from database into a map.
	stepsMap := make(map[string]string)
	stepsQuery := `SELECT id, code FROM steps WHERE deleted_at IS NULL`

	stepRows, err := s.db.QueryContext(ctx, stepsQuery)
	if err != nil {
		return fmt.Errorf("failed to load steps: %w", err)
	}
	defer stepRows.Close()

	for stepRows.Next() {
		var id, code string
		if err := stepRows.Scan(&id, &code); err != nil {
			return fmt.Errorf("failed to scan step: %w", err)
		}
		stepsMap[code] = id
	}

	// Group records by step code (in-memory, no DB queries in the loop).
	type stepData struct {
		headerTypeCodes []string
		recordCodes     []string
	}
	stepRecordsData := make(map[string]*stepData)

	var lastStepCode string
	var lastHeaderTypeCode string

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}

		stepCode := row[0]
		headerTypeCode := row[1]
		recordCode := row[2]

		if stepCode != "" {
			lastStepCode = stepCode
		} else {
			stepCode = lastStepCode
		}

		if headerTypeCode != "" {
			lastHeaderTypeCode = headerTypeCode
		} else {
			headerTypeCode = lastHeaderTypeCode
		}

		if stepCode == "" || headerTypeCode == "" || recordCode == "" {
			continue
		}

		sd, exists := stepRecordsData[stepCode]
		if !exists {
			sd = &stepData{}
			stepRecordsData[stepCode] = sd
		}

		if len(sd.headerTypeCodes) == 0 {
			headerTypes := strings.Split(headerTypeCode, ",")
			for _, ht := range headerTypes {
				trimmedHT := strings.TrimSpace(ht)
				if trimmedHT != "" {
					sd.headerTypeCodes = append(sd.headerTypeCodes, trimmedHT)
				}
			}
		}

		sd.recordCodes = append(sd.recordCodes, recordCode)
	}

	// Begin a single transaction for all inserts.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	processedSteps := 0
	createdStepHeaderTypes := 0
	createdStepRecords := 0
	skippedRows := 0

	for stepCode, data := range stepRecordsData {
		stepId, stepExists := stepsMap[stepCode]
		if !stepExists {
			skippedRows++
			continue
		}

		processedSteps++

		// Create step_header_types associations.
		for _, headerTypeCode := range data.headerTypeCodes {
			headerTypeId, headerTypeExists := headerTypesMap[headerTypeCode]
			if !headerTypeExists {
				skippedRows++
				continue
			}

			insertHeaderTypeQuery := `
				INSERT INTO step_header_types (step_id, header_type_id)
				VALUES ($1, $2)
				ON CONFLICT (step_id, header_type_id) DO NOTHING
			`

			_, err := tx.ExecContext(ctx, insertHeaderTypeQuery, stepId, headerTypeId)
			if err != nil {
				return fmt.Errorf("failed to create step_header_type for step %s and header type %s: %w", stepCode, headerTypeCode, err)
			}

			createdStepHeaderTypes++
		}

		// Create step_records associations.
		for _, recordCode := range data.recordCodes {
			recordId, recordExists := recordsMap[recordCode]
			if !recordExists {
				skippedRows++
				continue
			}

			insertQuery := `
				INSERT INTO step_records (step_id, record_id)
				VALUES ($1, $2)
				ON CONFLICT (step_id, record_id) DO NOTHING
			`

			_, err := tx.ExecContext(ctx, insertQuery, stepId, recordId)
			if err != nil {
				return fmt.Errorf("failed to create step_record for step %s and record %s: %w", stepCode, recordCode, err)
			}

			createdStepRecords++
		}

		// Silence per-step logging to keep output minimal.
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit step associations: %w", err)
	}

	return nil
}
