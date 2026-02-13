package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lantoniomiranda/shitreader/internal/app"
)

type fileConfig struct {
	path  string
	sheet string
}

type task struct {
	name string
	run  func() error
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

	tasks := make([]task, 0, len(files)+3)
	for _, file := range files {
		file := file
		tasks = append(tasks, task{
			name: fmt.Sprintf("Import %s", filepath.Base(file.path)),
			run: func() error {
				return app.ReaderService.Read(file.path, file.sheet)
			},
		})
	}

	tasks = append(tasks,
		task{
			name: "Inspect processo-passos.xlsx",
			run: func() error {
				return app.ReaderService.ReadProcessSteps("files/processo-passos.xlsx", "Data")
			},
		},
		task{
			name: "Associate fields with records",
			run: func() error {
				return app.AssociationService.Associate()
			},
		},
		task{
			name: "Associate record types",
			run: func() error {
				return app.AssociationService.AssociateRecordTypes("files/record-types.xlsx", "Data")
			},
		},
		task{
			name: "Associate steps",
			run: func() error {
				return app.AssociationService.AssociateSteps("files/passo-registos.xlsx", "Data")
			},
		},
	)

	totalTasks := len(tasks)
	start := time.Now()

	renderProgress(0, totalTasks, "Starting...", start)

	for i, task := range tasks {
		renderProgress(i, totalTasks, task.name, start)
		if err := task.run(); err != nil {
			log.Fatalf("Task %q failed: %v", task.name, err)
		}
	}

	renderProgress(totalTasks, totalTasks, "Completed\n", start)
	fmt.Printf("\nAll tasks finished in %s\n", time.Since(start).Round(time.Millisecond))
}

func renderProgress(completed, total int, current string, start time.Time) {
	if total == 0 {
		return
	}
	barWidth := 40
	percent := float64(completed) / float64(total)
	filled := int(percent * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Printf("\r[%s] %2.0f%% %d/%d | Elapsed: %s | Current: %s",
		bar,
		percent*100,
		completed,
		total,
		elapsed,
		current,
	)
}
