package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"workspace-onboarding-service/internal/models"
	"workspace-onboarding-service/internal/repository"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)


type UserService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

func NewUserService(userRepo repository.UserRepository, orgRepo repository.OrganizationRepository) *UserService {
	return &UserService{userRepo: userRepo, orgRepo: orgRepo}
}

func (s *UserService) Create(ctx context.Context, in models.CreateUserInput) (*models.User, error) {
	if err := validateCreateUser(in); err != nil {
		return nil, err
	}

	// Business rule from spec: "Organization must exist before creating a user."
	if _, err := s.orgRepo.GetByID(ctx, in.OrganizationID); err != nil {
		if err == models.ErrNotFound {
			return nil, models.ErrOrgNotFound
		}
		return nil, err
	}

	u := &models.User{
		OrganizationID: in.OrganizationID,
		Name:           strings.TrimSpace(in.Name),
		Email:          strings.ToLower(strings.TrimSpace(in.Email)),
		Role:           in.Role,
	}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Get(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, orgID *int64, limit, offset int) ([]models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // sane default + hard ceiling protects the DB from limit=1000000
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.List(ctx, orgID, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id int64, in models.UpdateUserInput) (*models.User, error) {
	existing, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		if strings.TrimSpace(*in.Name) == "" {
			return nil, fmt.Errorf("%w: name cannot be empty", models.ErrValidation)
		}
		existing.Name = strings.TrimSpace(*in.Name)
	}
	if in.Email != nil {
		if !emailRegex.MatchString(*in.Email) {
			return nil, fmt.Errorf("%w: invalid email", models.ErrValidation)
		}
		existing.Email = strings.ToLower(strings.TrimSpace(*in.Email))
	}
	if in.Role != nil {
		if !in.Role.Valid() {
			return nil, fmt.Errorf("%w: role must be admin or member", models.ErrValidation)
		}
		existing.Role = *in.Role
	}

	if err := s.userRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

func validateCreateUser(in models.CreateUserInput) error {
	if in.OrganizationID <= 0 {
		return fmt.Errorf("%w: organization_id is required", models.ErrValidation)
	}
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name is required", models.ErrValidation)
	}
	if !emailRegex.MatchString(in.Email) {
		return fmt.Errorf("%w: valid email is required", models.ErrValidation)
	}
	if !in.Role.Valid() {
		return fmt.Errorf("%w: role must be admin or member", models.ErrValidation)
	}
	return nil
}