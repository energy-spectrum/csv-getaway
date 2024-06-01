package entity

import (
	"csv-analyzer-api/internal/value"
)

type Filters struct {
	WorkYears int `json:"workYears"`
}

type Options struct {
	Priority float64 `json:"priority"` // from 0 to 1
}

type Skill struct {
	Name     string   `db:"name" json:"name"`
	Options  Options  `db:"options" json:"options"`
	Wordings []string `db:"wordings" json:"wordings"`
}

type Template struct {
	ID      value.TemplateID `db:"id" json:"id"`
	Name    string           `db:"name" json:"name"`
	Filters Filters          `db:"filters" json:"filters"`
	Skills  []Skill          `db:"skills" json:"skills"`
}
