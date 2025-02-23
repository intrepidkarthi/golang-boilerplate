package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for the Message model
func (Message) TableName() string {
	return "messages"
}

// BeforeUpdate hook to update the updated_at timestamp
func (m *Message) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
