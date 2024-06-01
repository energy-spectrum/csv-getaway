package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"csv-analyzer-api/internal/entity"

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

	if filters != nil {
		err = json.Unmarshal(filters, &i.Filters)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	if skills != nil {
		err = json.Unmarshal(skills, &i.Skills)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	return &i, nil
}

func (c Template) Delete(ctx context.Context, id value.TemplateID) error {
	const sqlQuery = `
		DELETE FROM
			public.templates
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
	row := c.db.QueryRowContext(
		ctx,
		sqlQuery,
		id,
	)

	var template entity.Template
	var filters, skills []byte
	err := row.Scan(
		&template.ID,
		&template.Name,
		&filters,
		&skills,
	)
	if err != nil {
		return nil, fmt.Errorf("db.Scan: %w", err)
	}

	if filters != nil {
		err = json.Unmarshal(filters, &template.Filters)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	if skills != nil {
		err = json.Unmarshal(skills, &template.Skills)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	return &template, nil
}

func (c Template) Get(ctx context.Context) ([]entity.Template, error) {
	sqlQuery := `
		SELECT *
		FROM
			public.templates
	`

	rows, err := c.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var templates []entity.Template
	for rows.Next() {
		var i entity.Template
		var filters, skills []byte
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&filters,
			&skills,
		); err != nil {
			return nil, err
		}

		if filters != nil {
			err = json.Unmarshal(filters, &i.Filters)
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %w", err)
			}
		}

		if skills != nil {
			err = json.Unmarshal(skills, &i.Skills)
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %w", err)
			}
		}

		templates = append(templates, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

func (c Template) Update(ctx context.Context, template *entity.Template) error {
	const sqlQuery = `
		UPDATE public.templates
		SET
			filters = COALESCE($1, filters),
			skills = COALESCE($2, skills),
			name = COALESCE($3, name)
		WHERE
			id = $4;
	`

	if _, err := c.db.ExecContext(
		ctx,
		sqlQuery,
		template.Filters,
		template.Skills,
		template.Name,
		template.ID,
	); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	// var i entity.Template
	// err := row.Scan(
	// 	&i.ID,
	// 	&i.Name,
	// 	&i.Filters,
	// 	&i.Skills,
	// )
	// if err != nil {
	// 	return fmt.Errorf("db.Scan: %w", err)
	// }

	return nil
}
