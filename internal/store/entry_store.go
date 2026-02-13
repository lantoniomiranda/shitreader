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

	// Normalization caches
	tableVersionCache map[string]string // "table_code|version" -> id
	catalogCache      map[string]string // "slug" -> id
}

func NewPostgresEntryStore(db *sql.DB) *PostgresEntryStore {
	return &PostgresEntryStore{
		db:                db,
		tableVersionCache: make(map[string]string),
		catalogCache:      make(map[string]string),
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
	case "steps", "records", "fields":
		return s.batchInsertStructural(ctx, tx, entries, tableName)
	default:
		return s.batchInsertCatalog(ctx, tx, entries, tableName)
	}
}

// Helpers for caches
func (s *PostgresEntryStore) getTableVersionID(ctx context.Context, tx *sql.Tx, tableCode, version string) (string, error) {
	key := tableCode + "|" + version
	if id, ok := s.tableVersionCache[key]; ok {
		return id, nil
	}

	var id string
	// Try to find existing
	err := tx.QueryRowContext(ctx, `SELECT id FROM table_versions WHERE table_code = $1 AND version = $2`, tableCode, version).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert new
		err = tx.QueryRowContext(ctx, `
			INSERT INTO table_versions (table_code, version) VALUES ($1, $2)
			ON CONFLICT (table_code, version) DO UPDATE SET updated_at = NOW()
			RETURNING id
		`, tableCode, version).Scan(&id)
	}
	if err != nil {
		return "", fmt.Errorf("failed to resolve table version for %s %s: %w", tableCode, version, err)
	}

	s.tableVersionCache[key] = id
	return id, nil
}

func (s *PostgresEntryStore) getCatalogID(ctx context.Context, tx *sql.Tx, slug string) (string, error) {
	if id, ok := s.catalogCache[slug]; ok {
		return id, nil
	}

	var id string
	err := tx.QueryRowContext(ctx, `SELECT id FROM catalogs WHERE slug = $1`, slug).Scan(&id)
	if err == sql.ErrNoRows {
		// Auto-create catalog
		name := strings.ReplaceAll(slug, "_", " ")
		name = strings.Title(name)
		err = tx.QueryRowContext(ctx, `
			INSERT INTO catalogs (slug, name) VALUES ($1, $2)
			ON CONFLICT (slug) DO UPDATE SET updated_at = NOW()
			RETURNING id
		`, slug, name).Scan(&id)
	}
	if err != nil {
		return "", fmt.Errorf("failed to resolve catalog for %s: %w", slug, err)
	}

	s.catalogCache[slug] = id
	return id, nil
}

// batchInsertStructural handles steps, records, fields (using table_version_id)
func (s *PostgresEntryStore) batchInsertStructural(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch to avoid "ON CONFLICT ... cannot affect row a second time"
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

	first := entries[0]
	tvId, err := s.getTableVersionID(ctx, tx, first.Table, first.Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, code, description)"
	colsPerRow := 3

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
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3))
			args = append(args, tvId, e.Code, e.Description)
		}

		// Use ON CONFLICT DO UPDATE to handle re-runs or duplicate rows in source
		conflictTarget := "(table_version_id, code)"
		updateSet := "SET description = EXCLUDED.description, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

// batchInsertCatalog handles all lookup tables (using catalog_values)
func (s *PostgresEntryStore) batchInsertCatalog(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

	catalogId, err := s.getCatalogID(ctx, tx, tableName)
	if err != nil {
		return err
	}

	first := entries[0]
	tvId, err := s.getTableVersionID(ctx, tx, first.Table, first.Version)
	if err != nil {
		return err
	}

	cols := "(catalog_id, table_version_id, code, description)"
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
			args = append(args, catalogId, tvId, e.Code, e.Description)
		}

		conflictTarget := "(catalog_id, table_version_id, code)"
		updateSet := "SET description = EXCLUDED.description, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO catalog_values %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into catalog_values (slug=%s): %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertCountries(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

	tvId, err := s.getTableVersionID(ctx, tx, entries[0].Table, entries[0].Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, code, name)"
	colsPerRow := 3

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
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3))
			args = append(args, tvId, e.Code, e.Name)
		}

		conflictTarget := "(table_version_id, code)"
		updateSet := "SET name = EXCLUDED.name, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertDistricts(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

	if s.countryPTId == "" {
		err := tx.QueryRowContext(ctx, `SELECT id FROM countries WHERE code = 'PT' AND deleted_at IS NULL`).Scan(&s.countryPTId)
		if err != nil {
			return fmt.Errorf("querying country PT: %w", err)
		}
	}

	tvId, err := s.getTableVersionID(ctx, tx, entries[0].Table, entries[0].Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, code, name, country_id)"
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
			args = append(args, tvId, e.Code, e.Name, s.countryPTId)
		}

		conflictTarget := "(table_version_id, code)"
		updateSet := "SET name = EXCLUDED.name, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertMunicipalities(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

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

	tvId, err := s.getTableVersionID(ctx, tx, entries[0].Table, entries[0].Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, code, name, district_id)"
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
			prefix := e.Code
			if len(e.Code) >= 2 {
				prefix = e.Code[:2]
			}
			districtId, ok := s.districtCache[prefix]
			if !ok {
				return fmt.Errorf("district not found for municipality code %s (prefix %s)", e.Code, prefix)
			}

			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4))
			args = append(args, tvId, e.Code, e.Name, districtId)
		}

		conflictTarget := "(table_version_id, code)"
		updateSet := "SET name = EXCLUDED.name, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertParishes(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.Code] {
			seen[e.Code] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

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

	tvId, err := s.getTableVersionID(ctx, tx, entries[0].Table, entries[0].Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, code, name, municipality_id)"
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
			prefix := e.Code
			if len(e.Code) >= 4 {
				prefix = e.Code[:4]
			}
			municipalId, ok := s.municipalCache[prefix]
			if !ok {
				return fmt.Errorf("municipality not found for parish code %s (prefix %s)", e.Code, prefix)
			}

			base := i * colsPerRow
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4))
			args = append(args, tvId, e.Code, e.Name, municipalId)
		}

		conflictTarget := "(table_version_id, code)"
		updateSet := "SET name = EXCLUDED.name, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}

func (s *PostgresEntryStore) batchInsertINEZones(ctx context.Context, tx *sql.Tx, entries []types.Entry, tableName string) error {
	// Deduplicate entries within the batch (using ZoneCode)
	seen := make(map[string]bool)
	uniqueEntries := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		if !seen[e.ZoneCode] {
			seen[e.ZoneCode] = true
			uniqueEntries = append(uniqueEntries, e)
		}
	}
	entries = uniqueEntries

	if len(entries) == 0 {
		return nil
	}

	tvId, err := s.getTableVersionID(ctx, tx, entries[0].Table, entries[0].Version)
	if err != nil {
		return err
	}

	cols := "(table_version_id, zone_code, zone_name, zone_name_formatted, ine_municipality_code)"
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
			args = append(args, tvId, e.ZoneCode, e.ZoneName, e.ZoneNameFormatted, e.INEMunicipalityCode)
		}

		conflictTarget := "(table_version_id, zone_code)"
		updateSet := "SET zone_name = EXCLUDED.zone_name, zone_name_formatted = EXCLUDED.zone_name_formatted, updated_at = NOW()"

		query := fmt.Sprintf("INSERT INTO %s %s VALUES %s ON CONFLICT %s DO UPDATE %s",
			tableName, cols, strings.Join(placeholders, ", "), conflictTarget, updateSet)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("batch insert into %s: %w", tableName, err)
		}
	}
	return nil
}
