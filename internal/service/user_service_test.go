package service_test

import (
	"context"
	"testing"

	"workspace-onboarding-service/internal/models"
	"workspace-onboarding-service/internal/service"
)



type fakeUserRepo struct {
	users  map[int64]*models.User
	nextID int64
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[int64]*models.User{}, nextID: 1}
}

func (f *fakeUserRepo) Create(_ context.Context, u *models.User) error {
	for _, existing := range f.users {
		if existing.Email == u.Email {
			return models.ErrDuplicate
		}
	}
	u.ID = f.nextID
	f.nextID++
	f.users[u.ID] = u
	return nil
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int64) (*models.User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, models.ErrNotFound
}

func (f *fakeUserRepo) List(_ context.Context, orgID *int64, _, _ int) ([]models.User, error) {
	var out []models.User
	for _, u := range f.users {
		if orgID == nil || u.OrganizationID == *orgID {
			out = append(out, *u)
		}
	}
	return out, nil
}

func (f *fakeUserRepo) Update(_ context.Context, u *models.User) error {
	if _, ok := f.users[u.ID]; !ok {
		return models.ErrNotFound
	}
	f.users[u.ID] = u
	return nil
}

func (f *fakeUserRepo) Delete(_ context.Context, id int64) error {
	if _, ok := f.users[id]; !ok {
		return models.ErrNotFound
	}
	delete(f.users, id)
	return nil
}


type fakeOrgRepo struct {
	orgs map[int64]*models.Organization
}

func newFakeOrgRepo(existingIDs ...int64) *fakeOrgRepo {
	orgs := map[int64]*models.Organization{}
	for _, id := range existingIDs {
		orgs[id] = &models.Organization{ID: id, Name: "Acme", Domain: "acme.com"}
	}
	return &fakeOrgRepo{orgs: orgs}
}

func (f *fakeOrgRepo) GetByID(_ context.Context, id int64) (*models.Organization, error) {
	if o, ok := f.orgs[id]; ok {
		return o, nil
	}
	return nil, models.ErrNotFound
}

func TestUserService_Create_Success(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	u, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 1, Name: "John Doe", Email: "john@acme.com", Role: models.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.ID == 0 {
		t.Fatalf("expected user to have an ID assigned")
	}
	if u.Email != "john@acme.com" {
		t.Fatalf("expected email to be normalized, got %q", u.Email)
	}
}

func TestUserService_Create_OrganizationMustExist(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo()) // no orgs exist

	_, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 99, Name: "John Doe", Email: "john@acme.com", Role: models.RoleMember,
	})
	if err != models.ErrOrgNotFound {
		t.Fatalf("expected ErrOrgNotFound, got %v", err)
	}
}

func TestUserService_Create_InvalidEmail(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	_, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 1, Name: "John Doe", Email: "not-an-email", Role: models.RoleMember,
	})
	if err == nil {
		t.Fatalf("expected a validation error, got nil")
	}
}

func TestUserService_Create_MissingName(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	_, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 1, Name: "   ", Email: "john@acme.com", Role: models.RoleMember,
	})
	if err == nil {
		t.Fatalf("expected a validation error for blank name, got nil")
	}
}

func TestUserService_Create_InvalidRole(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	_, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 1, Name: "John Doe", Email: "john@acme.com", Role: "superuser",
	})
	if err == nil {
		t.Fatalf("expected a validation error for invalid role, got nil")
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	input := models.CreateUserInput{OrganizationID: 1, Name: "John", Email: "john@acme.com", Role: models.RoleMember}
	if _, err := svc.Create(context.Background(), input); err != nil {
		t.Fatalf("first create should succeed, got %v", err)
	}
	if _, err := svc.Create(context.Background(), input); err != models.ErrDuplicate {
		t.Fatalf("expected ErrDuplicate on second create, got %v", err)
	}
}

func TestUserService_Update_NotFound(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	newName := "New Name"
	_, err := svc.Update(context.Background(), 12345, models.UpdateUserInput{Name: &newName})
	if err != models.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserService_Delete(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepo(), newFakeOrgRepo(1))

	u, err := svc.Create(context.Background(), models.CreateUserInput{
		OrganizationID: 1, Name: "John", Email: "john@acme.com", Role: models.RoleMember,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}
	if err := svc.Delete(context.Background(), u.ID); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}
	if _, err := svc.Get(context.Background(), u.ID); err != models.ErrNotFound {
		t.Fatalf("expected user to be gone after delete, got %v", err)
	}
}