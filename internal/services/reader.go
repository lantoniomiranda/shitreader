package services

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (s *ReaderService) Read(filePath string, sheetName string) error {
	fmt.Printf("\n📁 Processando arquivo: %s\n", filePath)
	
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

	totalRows := 0
	for _, row := range rows {
		if len(row) > 2 && row[2] != "" {
			totalRows++
		}
	}

	fmt.Printf("📊 Processando %d linhas de dados...\n\n", totalRows)

	processedRows := 0
	startTime := time.Now()
	var tableName string

	for _, row := range rows {
		isHeaderRow := len(row) >= 2 && (len(row) == 2 || (len(row) > 2 && row[2] == ""))
		
		if isHeaderRow {
			if t, ok := types.TableCodeMap[row[0]]; ok {
				tableName = t
			} else {
				tableName = "OTHER"
			}
			continue
		}

		if tableName == "" || tableName == "OTHER" {
			continue
		}

		entry := parseRow(row, tableName)

		if err := s.entryStore.Save(ctx, &entry, tableName); err != nil {
			return fmt.Errorf("error saving entry: %w", err)
		}

		processedRows++
		if processedRows%10 == 0 || processedRows == totalRows {
			s.printProgress(processedRows, totalRows, startTime, tableName)
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n\n✅ Processamento concluído em %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("📈 Total de registos processados: %d\n", processedRows)
	fmt.Printf("⚡ Velocidade média: %.2f registos/segundo\n", float64(processedRows)/elapsed.Seconds())

	return nil
}

func (s *ReaderService) printProgress(current, total int, startTime time.Time, tableName string) {
	percent := float64(current) / float64(total) * 100
	barWidth := 50
	filled := int(float64(barWidth) * float64(current) / float64(total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	elapsed := time.Since(startTime)
	rate := float64(current) / elapsed.Seconds()
	remaining := time.Duration(float64(total-current) / rate * float64(time.Second))

	fmt.Printf("\r[%s] %.1f%% | %d/%d | %s | ETA: %s | Tabela: %-20s",
		bar,
		percent,
		current,
		total,
		elapsed.Round(time.Millisecond),
		remaining.Round(time.Second),
		tableName,
	)
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
