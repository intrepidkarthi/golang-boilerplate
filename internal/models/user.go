package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Roles     []Role     `json:"roles,omitempty" db:"-"`
}

type Role struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
	Permissions []Permission `json:"permissions,omitempty" db:"-"`
}

type Permission struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Service     string     `json:"service" db:"service"`      // e.g. "billing"
	Resource    string     `json:"resource" db:"resource"`    // e.g. "invoice"
	Action      string     `json:"action" db:"action"`        // e.g. "read"
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
