package types

import "time"

type Entry struct {
	Id                  string
	Table               string
	Version             string
	Code                string
	Description         string
	Name                string
	ZoneCode            string
	ZoneName            string
	ZoneNameFormatted   string
	INEMunicipalityCode string
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
}
