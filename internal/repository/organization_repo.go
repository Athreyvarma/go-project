package repository

import (
	"context"

	"workspace-onboarding-service/internal/models"
)

type OrganizationRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Organization, error)
}