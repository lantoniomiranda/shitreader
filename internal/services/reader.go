package services

import (
	"context"
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

	// Begin a single transaction for the entire file import.
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
			// Table boundary: flush accumulated entries before switching.
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

		// Flush when the batch threshold is reached.
		if len(pendingEntries) >= flushThreshold {
			if err := s.entryStore.SaveBatch(ctx, tx, pendingEntries, pendingTable); err != nil {
				return fmt.Errorf("error saving batch: %w", err)
			}
			processedRows += len(pendingEntries)
			pendingEntries = pendingEntries[:0]
		}
	}

	// Flush any remaining entries.
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
