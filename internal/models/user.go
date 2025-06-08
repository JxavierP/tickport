package models

import (
	"time"

	"github.com/google/uuid"
)


type Role string

const (
	AdminRole Role = "admin"
	AgentRole Role = "agent"
)

func (r Role) IsValid() bool {
	return r == AdminRole || r == AgentRole
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // Hashed Password
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
