package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"csv-analyzer-api/internal/entity"
	errTemp "csv-analyzer-api/pkg/error_templates"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/value"
)

const (
	sqlUserFields = `
	id,
	email,
	name,
	hashed_password,
	role
`
)

type User struct {
	db  *sqlx.DB
	cfg *config.Configuration
}

func NewUser(db *sqlx.DB, cfg *config.Configuration) User {
	return User{
		db:  db,
		cfg: cfg,
	}
}

func (c User) Create(ctx context.Context, user entity.User) error {
	const sqlQuery = `
		INSERT INTO public.users (
			email,
			name,
			hashed_password,
			role
		)
		VALUES ($1,$2,$3,$4)
	`

	if _, err := c.db.ExecContext(
		ctx,
		sqlQuery,
		user.Email,
		user.Name,
		user.HashedPassword,
		user.Role,
	); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}

func (c User) Delete(ctx context.Context, userID value.UserID) error {
	const sqlQuery = `
		DELETE FROM
			public.user
		WHERE
			call_list_id = $1
	`

	if _, err := c.db.ExecContext(
		ctx,
		sqlQuery,
		userID,
	); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func (c User) GetByID(ctx context.Context, id value.UserID) (*entity.User, error) {
	sqlQuery := `
		SELECT
			` + sqlUserFields + `
		FROM
			public.users
		WHERE
			id = $1
	`
	var user entity.User
	if err := c.db.GetContext(
		ctx,
		&user,
		sqlQuery,
		id,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTemp.New("user not found", http.StatusNotFound)
		}

		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &user, nil
}

func (c User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sqlQuery := `
		SELECT
			` + sqlUserFields + `
		FROM
			public.users
		WHERE
			email = $1
	`
	var user entity.User
	if err := c.db.GetContext(
		ctx,
		&user,
		sqlQuery,
		email,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTemp.New("user not found", http.StatusNotFound)
		}

		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &user, nil
}

func (c User) Update(
	ctx context.Context,
	user *entity.User,
) error {
	const sqlQuery = `
		UPDATE
			public.user
		SET
			role = $1
		WHERE
		    id = $2
	`

	row := c.db.QueryRowContext(
		ctx,
		sqlQuery,
		user.Role,
		user.ID,
	)

	var i entity.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.HashedPassword,
	)
	if err != nil {
		return fmt.Errorf("db.Scan: %w", err)
	}

	return err
}
