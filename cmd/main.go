package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lantoniomiranda/shitreader/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("WARNING: .env not loaded: %v", err)
	}

	app, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err) // This will show the actual error
	}
	defer app.DB.Close()

	// Process main data tables
	if err := app.ReaderService.Read("files/tabelas-dados.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read tabelas-dados.xlsx: %v", err)
	}

	// Process CAE Rev.4 data
	if err := app.ReaderService.Read("files/cae.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read cae.xlsx: %v", err)
	}

	// Process geographic data - countries (paises)
	if err := app.ReaderService.Read("files/distritos.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read distritos.xlsx: %v", err)
	}

	// Process geographic data - countries (paises)
	if err := app.ReaderService.Read("files/paises.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read paises.xlsx: %v", err)
	}

	// Process geographic data - municipalities (concelhos)
	if err := app.ReaderService.Read("files/concelhos.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read concelhos.xlsx: %v", err)
	}

	// Process geographic data - parishes (freguesias)
	if err := app.ReaderService.Read("files/freguesias.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read freguesias.xlsx: %v", err)
	}

	// Process geographic data - INE zones
	if err := app.ReaderService.Read("files/ine-zonas.xlsx", "Data"); err != nil {
		log.Fatalf("Failed to read ine-zonas.xlsx: %v", err)
	}

	log.Println("✅ All data imported successfully!")
}
