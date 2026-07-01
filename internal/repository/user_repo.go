package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"


	"workspace-onboarding-service/internal/models"
)


type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	List(ctx context.Context, orgID *int64, limit, offset int) ([]models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	const query = `
		INSERT INTO users (organization_id, name, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query, u.OrganizationID, u.Name, u.Email, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return mapPgError(err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const query = `
		SELECT id, organization_id, name, email, role, created_at, updated_at
		FROM users WHERE id = $1`

	var u models.User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&u.ID, &u.OrganizationID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}


func (r *userRepository) List(ctx context.Context, orgID *int64, limit, offset int) ([]models.User, error) {
	var rows pgx.Rows
	var err error

	if orgID != nil {
		rows, err = r.pool.Query(ctx, `
			SELECT id, organization_id, name, email, role, created_at, updated_at
			FROM users WHERE organization_id = $1
			ORDER BY id LIMIT $2 OFFSET $3`, *orgID, limit, offset)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, organization_id, name, email, role, created_at, updated_at
			FROM users ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0, limit)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.OrganizationID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) Update(ctx context.Context, u *models.User) error {
	const query = `
		UPDATE users SET name = $1, email = $2, role = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query, u.Name, u.Email, u.Role, u.ID).Scan(&u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNotFound
		}
		return mapPgError(err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}


func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation (duplicate email)
			return models.ErrDuplicate
		case "23503": // foreign_key_violation (organization_id doesn't exist)
			return models.ErrOrgNotFound
		}
	}
	return err
}