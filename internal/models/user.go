package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
	Active   bool      `gorm:"default:true"`
	Roles    []Role    `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"uniqueIndex;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Service    string `gorm:"index;not null"` // e.g. "billing"
	Resource   string `gorm:"index;not null"` // e.g. "invoice"
	Action     string `gorm:"index;not null"` // e.g. "read"
	Description string
}
