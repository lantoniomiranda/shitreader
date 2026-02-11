package types

import "time"

type Entry struct {
	Id          string
	Table       string
	Version     string
	Code        string
	Description string

	// Additional fields for geographic tables (countries, districts, municipalities, parishes)
	Name string

	// Additional fields for INE zones
	ZoneCode            string
	ZoneName            string
	ZoneNameFormatted   string
	INEMunicipalityCode string

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
