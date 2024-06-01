package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"csv-analyzer-api/internal/entity"
	errTemp "csv-analyzer-api/pkg/error_templates"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/value"
)

type Template struct {
	db  *sqlx.DB
	cfg *config.Configuration
}

func NewTemplate(db *sqlx.DB, cfg *config.Configuration) Template {
	return Template{
		db:  db,
		cfg: cfg,
	}
}

func (c Template) Create(ctx context.Context, template *entity.Template) (*entity.Template, error) {
	const sqlQuery = `
		INSERT INTO public.templates (
			name,
			filters,
			skills
		)
		VALUES ($1,$2,$3)
		RETURNING id, name, filters, skills
	`

	row := c.db.QueryRowContext(
		ctx,
		sqlQuery,
		template.Name,
		template.Filters,
		template.Skills,
	)

	var i entity.Template
	var filters, skills []byte 
	err := row.Scan(
		&i.ID,
		&i.Name,
		&filters,
		&skills,
	)
	if err != nil {
		return nil, fmt.Errorf("db.Scan: %w", err)
	}

	err = json.Unmarshal(filters, &i.Filters)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	err = json.Unmarshal(skills, &i.Skills)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &i, nil
}

func (c Template) Delete(ctx context.Context, id value.TemplateID) error {
	const sqlQuery = `
		DELETE FROM
			public.template
		WHERE
			id = $1
	`

	if _, err := c.db.ExecContext(
		ctx,
		sqlQuery,
		id,
	); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func (c Template) GetByID(ctx context.Context, id value.TemplateID) (*entity.Template, error) {
	sqlQuery := `
		SELECT *
		FROM
			public.templates
		WHERE
			id = $1
	`
	var template entity.Template
	if err := c.db.GetContext(
		ctx,
		&template,
		sqlQuery,
		id,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTemp.New("template not found", http.StatusNotFound)
		}

		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &template, nil
}


func (c Template) Update(ctx context.Context, template *entity.Template) error {
	const sqlQuery = `
		UPDATE public.template
		SET
			filters = COALESCE($1, filters),
			skills = COALESCE($2, skills),
			name = COALESCE($3, name)
		WHERE
			id = $4;
	`

	row := c.db.QueryRowContext(
		ctx,
		sqlQuery,
		template.Filters,
		template.ID,
	)

	var i entity.Template
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Filters,
		&i.Skills,
	)
	if err != nil {
		return fmt.Errorf("db.Scan: %w", err)
	}

	return err
}
