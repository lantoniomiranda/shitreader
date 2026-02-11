package store

import (
	"context"
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
	Save(ctx context.Context, entry *types.Entry, tableName string) error
}

func (s *PostgresEntryStore) Save(ctx context.Context, entry *types.Entry, tableName string) error {
	var query string
	var err error

	switch tableName {
	case "countries":
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, name, code)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
			query,
			entry.Table,
			entry.Version,
			entry.Name,
			entry.Code,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	case "districts":
		var countryId string
		queryCountry := `
		SELECT id FROM countries WHERE code = 'PT' AND deleted_at IS NULL
		`

		err = s.db.QueryRowContext(ctx, queryCountry).Scan(&countryId)
		if err != nil {
			return fmt.Errorf("Error querying country: %w", err)
		}

		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, name, code, country_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
			query,
			entry.Table,
			entry.Version,
			entry.Name,
			entry.Code,
			countryId,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	case "municipalities":
		var districtId string
		queryDistrict := `
		SELECT id FROM districts WHERE code LIKE $1 AND deleted_at IS NULL
		`
		firstTwoDigits := ""
		if len(entry.Code) >= 2 {
			firstTwoDigits = entry.Code[:2]
		} else {
			firstTwoDigits = entry.Code
		}

		queryCondition := firstTwoDigits + "%"
		err = s.db.QueryRowContext(ctx, queryDistrict, queryCondition).Scan(&districtId)
		if err != nil {
			return fmt.Errorf("Error querying district: %w", err)
		}

		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, name, code, district_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
			query,
			entry.Table,
			entry.Version,
			entry.Name,
			entry.Code,
			districtId,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	case "parishes":
		var municipalityId string
		queryMunicipality := `
		SELECT id FROM municipalities WHERE code LIKE $1 AND deleted_at IS NULL
		`
		digits := ""
		if len(entry.Code) >= 4 {
			digits = entry.Code[:4]
		} else {
			digits = entry.Code
		}

		queryCondition := digits + "%"
		err = s.db.QueryRowContext(ctx, queryMunicipality, queryCondition).Scan(&municipalityId)
		if err != nil {
			return fmt.Errorf("Error querying municipality: %w", err)
		}

		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, name, code, municipality_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
			query,
			entry.Table,
			entry.Version,
			entry.Name,
			entry.Code,
			municipalityId,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)
	case "ine_zones":
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, zone_code, zone_name, zone_name_formatted, ine_municipality_code)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
			query,
			entry.Table,
			entry.Version,
			entry.ZoneCode,
			entry.ZoneName,
			entry.ZoneNameFormatted,
			entry.INEMunicipalityCode,
		).Scan(&entry.Id, &entry.CreatedAt, &entry.UpdatedAt)

	default:
		query = fmt.Sprintf(`
			INSERT INTO %s (table_code, version, code, description)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at 
		`, tableName)

		err = s.db.QueryRowContext(
			ctx,
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
