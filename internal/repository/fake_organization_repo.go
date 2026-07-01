package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"workspace-onboarding-service/internal/models"
)


type localOrganizationRepo struct {
	pool *pgxpool.Pool
}

func NewLocalOrganizationRepo(pool *pgxpool.Pool) OrganizationRepository {
	return &localOrganizationRepo{pool: pool}
}

func (r *localOrganizationRepo) GetByID(ctx context.Context, id int64) (*models.Organization, error) {
	const query = `SELECT id, name, domain, created_at, updated_at FROM organizations WHERE id = $1`

	var o models.Organization
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&o.ID, &o.Name, &o.Domain, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}