package store

import (
	"database/sql"
	"fmt"

	"github.com/lantoniomiranda/shitreader/internal/types"
)

type PostgresEntryStore struct {
	db *sql.DB
}

func NewPostgresEntryStore(db *sql.DB) *PostgresEntryStore {
	return &PostgresEntryStore{
		db: db,
	}
}

type EntryStore interface {
	Save(entry *types.Entry, tableName string) error
}

func (s *PostgresEntryStore) Save(entry *types.Entry, tableName string) error {
	var query string
	var err error

	switch tableName {
	case "countries", "districts", "municipalities", "parishes":
		// Geographic tables with 'name' field
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, name, code)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRow(
			query,
			entry.Table,
			entry.Version,
			entry.Name,
			entry.Code,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	case "ine_zones":
		// INE zones with completely different structure
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, zone_code, zone_name, zone_name_formatted, ine_municipality_code)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRow(
			query,
			entry.Table,
			entry.Version,
			entry.ZoneCode,
			entry.ZoneName,
			entry.ZoneNameFormatted,
			entry.INEMunicipalityCode,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	default:
		// Standard tables (including cae_rev4 and postal_codes)
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, code, description)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRow(
			query,
			entry.Table,
			entry.Version,
			entry.Code,
			entry.Description,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)
	}

	if err != nil {
		return fmt.Errorf("Error saving entry to %s: %w", tableName, err)
	}

	return nil
}
