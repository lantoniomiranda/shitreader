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

func (s *PostgresAssociationStore) AssociateRecordsFields(ctx context.Context) error {
	fmt.Printf("\n🔗 Associating fields with records...\n")

	selectFieldsQuery := `
		SELECT id, code, table_code, version FROM fields WHERE deleted_at IS NULL
	`

	rows, err := s.db.QueryContext(ctx, selectFieldsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	associatedCount := 0
	skippedCount := 0

	for rows.Next() {
		var id, code, tableCode, version string

		err := rows.Scan(&id, &code, &tableCode, &version)
		if err != nil {
			return err
		}

		recordRoot := code[:5] + "%"

		var recordId string
		getRecordIdQuery := `
			SELECT id FROM records 
			WHERE code LIKE $1 
			AND deleted_at IS NULL
			LIMIT 1
		`

		err = s.db.QueryRowContext(ctx, getRecordIdQuery, recordRoot).Scan(&recordId)
		if err != nil {
			if err == sql.ErrNoRows {
				skippedCount++
				continue
			}
			return fmt.Errorf("failed to query record for field %s (code: %s, table: %s, version: %s): %w", id, code, tableCode, version, err)
		}

		updateFieldQuery := `
			UPDATE fields SET record_id = $1 WHERE id = $2
		`

		_, err = s.db.ExecContext(ctx, updateFieldQuery, recordId, id)
		if err != nil {
			return fmt.Errorf("failed to update field %s with record_id %s: %w", id, recordId, err)
		}

		associatedCount++
		if associatedCount%100 == 0 {
			fmt.Printf("  ✓ %d fields processed...\n", associatedCount)
		}
	}

	fmt.Printf("✅ Completed: %d fields associated, %d skipped\n", associatedCount, skippedCount)

	return nil
}

func (s *PostgresAssociationStore) AssociateRecordsRecordTypes(ctx context.Context, filePath string, sheetName string) error {
	fmt.Printf("\n🔗 Associating records with record types...\n")

	// Open Excel file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Read rows from Excel
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	// Load all record types from database into a map
	recordTypesMap := make(map[string]string) // key: record_type_code, value: record_type_id
	recordTypesQuery := `SELECT id, code FROM record_types WHERE deleted_at IS NULL`

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

	fmt.Printf("  Loaded %d record types\n", len(recordTypesMap))

	associatedCount := 0
	skippedCount := 0

	// Iterate through Excel rows (skip header row)
	for i, row := range rows {
		// Skip header row
		if i == 0 {
			continue
		}

		// Skip rows that don't have at least 3 columns
		if len(row) < 3 {
			continue
		}

		recordCode := row[1]     // Column 1: Record code (e.g., R000000)
		recordTypeCode := row[2] // Column 2: Record type (e.g., 1, 2)

		// Skip empty rows
		if recordCode == "" || recordTypeCode == "" {
			continue
		}

		// Find the record_type_id from the map
		recordTypeId, exists := recordTypesMap[recordTypeCode]
		if !exists {
			skippedCount++
			continue
		}

		// Update the record with the record_type_id
		updateQuery := `
			UPDATE records 
			SET record_type_id = $1 
			WHERE code = $2 
			AND deleted_at IS NULL
		`

		result, err := s.db.ExecContext(ctx, updateQuery, recordTypeId, recordCode)
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
		if associatedCount%100 == 0 {
			fmt.Printf("  ✓ %d records processed...\n", associatedCount)
		}
	}

	fmt.Printf("✅ Completed: %d records associated, %d skipped\n", associatedCount, skippedCount)

	return nil
}

func (s *PostgresAssociationStore) AssociateStepsHeaderTypesAndRecords(ctx context.Context, filePath string, sheetName string) error {
	fmt.Printf("\n🔗 Associating steps with header types and records...\n")

	// Open Excel file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Read rows from Excel
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	// Load all header types from database into a map
	headerTypesMap := make(map[string]string) // key: header_type_code, value: header_type_id
	headerTypesQuery := `SELECT id, code FROM header_types WHERE deleted_at IS NULL`

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

	fmt.Printf("  Loaded %d header types\n", len(headerTypesMap))

	// Load all records from database into a map
	recordsMap := make(map[string]string) // key: record_code, value: record_id
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

	fmt.Printf("  Loaded %d records\n", len(recordsMap))

	// Load all steps from database into a map
	stepsMap := make(map[string]string) // key: step_code, value: step_id
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

	fmt.Printf("  Loaded %d steps\n", len(stepsMap))

	processedSteps := 0
	createdStepHeaderTypes := 0
	createdStepRecords := 0
	skippedRows := 0

	// Group records by step code
	stepRecordsData := make(map[string]struct {
		headerTypeCodes []string
		recordCodes     []string
	})

	// Track the last non-empty values for merged cells
	var lastStepCode string
	var lastHeaderTypeCode string

	// Iterate through Excel rows (skip header row)
	for i, row := range rows {
		// Skip header row
		if i == 0 {
			continue
		}

		// Skip rows that don't have at least 3 columns (need: passo, tipo_cab, registo)
		if len(row) < 3 {
			continue
		}

		stepCode := row[0]       // Column 0: "Passo" (e.g., P1100, P4120)
		headerTypeCode := row[1] // Column 1: "Tipo Cab." (e.g., I or I,O)
		recordCode := row[2]     // Column 2: "Registo" (e.g., R000000, R110000)

		// Handle merged cells: use the last non-empty value if current is empty
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

		// Skip if we still don't have values or record code is empty
		if stepCode == "" || headerTypeCode == "" || recordCode == "" {
			continue
		}

		// Store data grouped by step
		if _, exists := stepRecordsData[stepCode]; !exists {
			stepRecordsData[stepCode] = struct {
				headerTypeCodes []string
				recordCodes     []string
			}{
				headerTypeCodes: []string{},
				recordCodes:     []string{},
			}
		}

		data := stepRecordsData[stepCode]

		// Handle comma-separated header types
		if data.headerTypeCodes == nil || len(data.headerTypeCodes) == 0 {
			// Parse comma-separated header types
			headerTypes := strings.Split(headerTypeCode, ",")
			for _, ht := range headerTypes {
				trimmedHT := strings.TrimSpace(ht)
				if trimmedHT != "" {
					data.headerTypeCodes = append(data.headerTypeCodes, trimmedHT)
				}
			}
		}

		data.recordCodes = append(data.recordCodes, recordCode)
		stepRecordsData[stepCode] = data
	}

	// Process each step
	for stepCode, data := range stepRecordsData {
		// Get step_id
		stepId, stepExists := stepsMap[stepCode]
		if !stepExists {
			skippedRows++
			continue
		}

		processedSteps++

		// Create step_header_types associations for ALL header types
		for _, headerTypeCode := range data.headerTypeCodes {
			headerTypeId, headerTypeExists := headerTypesMap[headerTypeCode]
			if !headerTypeExists {
				skippedRows++
				continue
			}

			// Insert step_header_type association
			insertHeaderTypeQuery := `
				INSERT INTO step_header_types (step_id, header_type_id)
				VALUES ($1, $2)
				ON CONFLICT (step_id, header_type_id) DO NOTHING
			`

			_, err := s.db.ExecContext(ctx, insertHeaderTypeQuery, stepId, headerTypeId)
			if err != nil {
				return fmt.Errorf("failed to create step_header_type for step %s and header type %s: %w", stepCode, headerTypeCode, err)
			}

			createdStepHeaderTypes++
		}

		// Create step_records associations
		for _, recordCode := range data.recordCodes {
			recordId, recordExists := recordsMap[recordCode]
			if !recordExists {
				skippedRows++
				continue
			}

			// Insert step_record association
			insertQuery := `
				INSERT INTO step_records (step_id, record_id)
				VALUES ($1, $2)
				ON CONFLICT (step_id, record_id) DO NOTHING
			`

			_, err := s.db.ExecContext(ctx, insertQuery, stepId, recordId)
			if err != nil {
				return fmt.Errorf("failed to create step_record for step %s and record %s: %w", stepCode, recordCode, err)
			}

			createdStepRecords++
		}

		if processedSteps%10 == 0 {
			fmt.Printf("  ✓ %d steps processed...\n", processedSteps)
		}
	}

	fmt.Printf("✅ Completed: %d steps processed, %d step-header_type associations, %d step-record associations, %d rows skipped\n",
		processedSteps, createdStepHeaderTypes, createdStepRecords, skippedRows)

	return nil
}
