package models

import "time"

type Role string

const(
	RoleAdmin Role = "admin"
	RoleMember Role = "member"
)

func (r Role) Valid() bool{
	return r == RoleAdmin || r == RoleMember
}

type User struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	OrganizationID int64  `json:"organization_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           Role   `json:"role"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Role  *Role   `json:"role,omitempty"`
}