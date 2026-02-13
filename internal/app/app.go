package app

import (
	"database/sql"
	"fmt"

	"github.com/lantoniomiranda/shitreader/internal/services"
	"github.com/lantoniomiranda/shitreader/internal/store"
	"github.com/lantoniomiranda/shitreader/migrations"
)

type Application struct {
	ReaderService      *services.ReaderService
	AssociationService *services.AssociationService
	DB                 *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	if err := pgDb.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	entryStore := store.NewPostgresEntryStore(pgDb)
	associationStore := store.NewPostgresAssociationStore(pgDb)

	readerService := services.NewReaderService(entryStore)
	associationService := services.NewAssociationService(associationStore)

	return &Application{
		ReaderService:      readerService,
		AssociationService: associationService,
		DB:                 pgDb,
	}, nil
}
