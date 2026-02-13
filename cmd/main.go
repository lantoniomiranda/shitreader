package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lantoniomiranda/shitreader/internal/app"
)

type fileConfig struct {
	path  string
	sheet string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("WARNING: .env not loaded: %v", err)
	}

	app, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	defer app.DB.Close()

	files := []fileConfig{
		{"files/tabelas-dados.xlsx", "Data"},
		{"files/cae.xlsx", "Data"},
		{"files/paises.xlsx", "Data"},
		{"files/distritos.xlsx", "Data"},
		{"files/concelhos.xlsx", "Data"},
		{"files/freguesias.xlsx", "Data"},
		{"files/ine-zonas.xlsx", "Data"},
	}

	for _, file := range files {
		if err := app.ReaderService.Read(file.path, file.sheet); err != nil {
			log.Fatalf("Failed to read %s: %v", file.path, err)
		}
	}

	log.Println("✅ All data imported successfully!")

	log.Println("\n🔗 Starting associations...")

	err = app.AssociationService.Associate()
	if err != nil {
		log.Fatalf("Failed to associate fields with records: %v", err)
	}

	err = app.AssociationService.AssociateRecordTypes("files/record-types.xlsx", "Data")
	if err != nil {
		log.Fatalf("Failed to associate records with record types: %v", err)
	}

	err = app.AssociationService.AssociateSteps("files/passo-registos.xlsx", "Data")
	if err != nil {
		log.Fatalf("Failed to associate steps with header types and records: %v", err)
	}

	log.Println("\n✅ All associations completed successfully!")
}
