package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lantoniomiranda/shitreader/internal/types"
)

type PostgresEntryStore struct {
	db *sql.DB

	// Cached parent lookups, populated lazily.
	countryPTId    string
	districtCache  map[string]string // code prefix (2 digits) -> id
	municipalCache map[string]string // code prefix (4 digits) -> id
}

func NewPostgresEntryStore(db *sql.DB) *PostgresEntryStore {
	return &PostgresEntryStore{
		db: db,
	}
}

type EntryStore interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	SaveBatch(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error
}

const batchSize = 500

func (s *PostgresEntryStore) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *PostgresEntryStore) SaveBatch(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	if len(entries) == 0 {
		return nil
	}

	switch tableName {
	case "countries":
		return s.batchInsertCountries(ctx, tx, entries, tableName)
	case "districts":
		return s.batchInsertDistricts(ctx, tx, entries, tableName)
	case "municipalities":
		return s.batchInsertMunicipalities(ctx, tx, entries, tableName)
	case "parishes":
		return s.batchInsertParishes(ctx, tx, entries, tableName)
	case "ine_zones":
		return s.batchInsertINEZones(ctx, tx, entries, tableName)
	default:
		return s.batchInsertDefault(ctx, tx, entries, tableName)
	}
}

func (s *PostgresEntryStore) batchInsertCountries(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	cols := "(table_code, version, name, code)"
	colsPerRow := 4

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4))
			args = append(args, e.Table, e.Version, e.Name, e.Code)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertDistricts(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	if s.countryPTId == "" {
		err := tx.QueryRowContext(ctx, `SELECT id FROM countries WHERE code = 'PT' AND deleted_at IS NULL`).Scan(&s.countryPTId)
		if err != nil {
			return fmt.Errorf("querying country PT: %w", err)
		}
	}

	cols := "(table_code, version, name, code, country_id)"
	colsPerRow := 5

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5))
			args = append(args, e.Table, e.Version, e.Name, e.Code, s.countryPTId)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertMunicipalities(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	if s.districtCache == nil {
		s.districtCache = make(map[string]string)
		rows, err := tx.QueryContext(ctx, `SELECT id, code FROM districts WHERE deleted_at IS NULL`)
		if err != nil {
			return fmt.Errorf("loading districts cache: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var id, code string
			if err := rows.Scan(&id, &code); err != nil {
				return fmt.Errorf("scanning district: %w", err)
			}
			if len(code) >= 2 {
				s.districtCache[code[:2]] = id
			}
		}
	}

	cols := "(table_code, version, name, code, district_id)"
	colsPerRow := 5

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			prefix := e.Code
			if len(e.Code) >= 2 {
				prefix = e.Code[:2]
			}
			districtId, ok := s.districtCache[prefix]
			if !ok {
				return fmt.Errorf("district not found for municipality code %s (prefix %s)", e.Code, prefix)
			}

			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5))
			args = append(args, e.Table, e.Version, e.Name, e.Code, districtId)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertParishes(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	if s.municipalCache == nil {
		s.municipalCache = make(map[string]string)
		rows, err := tx.QueryContext(ctx, `SELECT id, code FROM municipalities WHERE deleted_at IS NULL`)
		if err != nil {
			return fmt.Errorf("loading municipalities cache: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var id, code string
			if err := rows.Scan(&id, &code); err != nil {
				return fmt.Errorf("scanning municipality: %w", err)
			}
			if len(code) >= 4 {
				s.municipalCache[code[:4]] = id
			}
		}
	}

	cols := "(table_code, version, name, code, municipality_id)"
	colsPerRow := 5

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			prefix := e.Code
			if len(e.Code) >= 4 {
				prefix = e.Code[:4]
			}
			municipalId, ok := s.municipalCache[prefix]
			if !ok {
				return fmt.Errorf("municipality not found for parish code %s (prefix %s)", e.Code, prefix)
			}

			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5))
			args = append(args, e.Table, e.Version, e.Name, e.Code, municipalId)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertINEZones(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	cols := "(table_code, version, zone_code, zone_name, zone_name_formatted, ine_municipality_code)"
	colsPerRow := 6

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5, base+6))
			args = append(args, e.Table, e.Version, e.ZoneCode, e.ZoneName, e.ZoneNameFormatted, e.INEMunicipalityCode)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertDefault(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	cols := "(table_code, version, code, description)"
	colsPerRow := 4

	for start := 0; start < len(entries); start += batchSize {
		end := start + batchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[start:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*colsPerRow)
		for i, e := range batch {
			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4))
			args = append(args, e.Table, e.Version, e.Code, e.Description)
		}

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, strings.Join(placeholders, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}
