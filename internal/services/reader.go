package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lantoniomiranda/shitreader/internal/store"
	"github.com/lantoniomiranda/shitreader/internal/types"
	"github.com/xuri/excelize/v2"
)

type ReaderService struct {
	entryStore store.EntryStore
}

func NewReaderService(entryStore store.EntryStore) *ReaderService {
	return &ReaderService{
		entryStore: entryStore,
	}
}

const flushThreshold = 500

func (s *ReaderService) Read(filePath string, sheetName string) error {
	ctx := context.Background()

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	tx, err := s.entryStore.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	processedRows := 0
	var tableName string
	var pendingEntries []types.Entry
	var pendingTable string

	for _, row := range rows {
		isHeaderRow := len(row) >= 2 && (len(row) == 2 || (len(row) > 2 && row[2] == ""))

		if isHeaderRow {
			if len(pendingEntries) > 0 {
				if err := s.entryStore.SaveBatch(ctx, tx, pendingEntries, pendingTable); err != nil {
					return fmt.Errorf("error saving batch: %w", err)
				}
				processedRows += len(pendingEntries)
				pendingEntries = pendingEntries[:0]
			}

			if t, ok := types.TableCodeMap[row[0]]; ok {
				tableName = t
			} else {
				tableName = "OTHER"
			}
			pendingTable = tableName
			continue
		}

		if tableName == "" || tableName == "OTHER" {
			continue
		}

		entry := parseRow(row, tableName)
		pendingEntries = append(pendingEntries, entry)

		if len(pendingEntries) >= flushThreshold {
			if err := s.entryStore.SaveBatch(ctx, tx, pendingEntries, pendingTable); err != nil {
				return fmt.Errorf("error saving batch: %w", err)
			}
			processedRows += len(pendingEntries)
			pendingEntries = pendingEntries[:0]
		}
	}

	if len(pendingEntries) > 0 {
		if err := s.entryStore.SaveBatch(ctx, tx, pendingEntries, pendingTable); err != nil {
			return fmt.Errorf("error saving batch: %w", err)
		}
		processedRows += len(pendingEntries)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s *ReaderService) ReadProcessSteps(filePath string, sheetName string) error {
	ctx := context.Background()

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	if len(rows) == 0 {
		return nil
	}

	var lastProcesso string

	processesMap := make(map[string][]string)

	for i, row := range rows {
		if i == 0 {
			continue
		}
		var processo, passo string
		if len(row) > 0 {
			processo = row[0]
		}
		if len(row) > 1 {
			passo = row[1]
		}

		if processo == "" {
			processo = lastProcesso
		} else {
			lastProcesso = processo
		}

		if passo == "" {
			continue
		}

		if processo == "" {
			continue
		}

		processesMap[processo] = append(processesMap[processo], passo)
	}

	tx, err := s.entryStore.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	for processCode, stepCodes := range processesMap {
		var description string
		err := tx.QueryRowContext(ctx, `
			SELECT COALESCE(cv.description, '')
			FROM catalog_values cv
			JOIN catalogs c ON cv.catalog_id = c.id
			WHERE c.slug = 'processes' AND cv.code = $1
			LIMIT 1
		`, processCode).Scan(&description)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error fetching process description %s: %w", processCode, err)
		}

		var processID string
		err = tx.QueryRowContext(ctx, `
			INSERT INTO processes (code, description)
			VALUES ($1, $2)
			ON CONFLICT (code) DO UPDATE SET description = EXCLUDED.description, updated_at = NOW()
			RETURNING id
		`, processCode, description).Scan(&processID)
		if err != nil {
			return fmt.Errorf("error saving process %s: %w", processCode, err)
		}

		for stepOrder, stepCode := range stepCodes {
			var stepID string
			err := tx.QueryRowContext(ctx, `
				SELECT id FROM steps WHERE code = $1 AND deleted_at IS NULL LIMIT 1
			`, stepCode).Scan(&stepID)
			if err != nil {
				continue
			}

			_, err = tx.ExecContext(ctx, `
				INSERT INTO process_steps (process_id, step_id, step_order)
				VALUES ($1, $2, $3)
				ON CONFLICT (process_id, step_id) DO UPDATE SET step_order = $3, updated_at = NOW()
			`, processID, stepID, stepOrder+1)
			if err != nil {
				return fmt.Errorf("error linking process %s to step %s: %w", processCode, stepCode, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func parseRow(row []string, tableName string) types.Entry {
	entry := types.Entry{}

	if len(row) > 0 {
		entry.Table = row[0]
	}
	if len(row) > 1 {
		entry.Version = row[1]
	}

	switch tableName {
	case "countries", "districts", "municipalities", "parishes":
		if len(row) > 2 {
			entry.Name = row[2]
		}
		if len(row) > 3 {
			entry.Code = row[3]
		}

	case "ine_zones":
		if len(row) > 2 {
			entry.ZoneCode = row[2]
		}
		if len(row) > 3 {
			entry.ZoneName = row[3]
		}
		if len(row) > 4 {
			entry.ZoneNameFormatted = row[4]
		}
		if len(row) > 5 {
			entry.INEMunicipalityCode = row[5]
		}

	default:
		if len(row) > 2 {
			entry.Code = row[2]
		}
		if len(row) > 3 {
			entry.Description = row[3]
		}
	}

	return entry
}
