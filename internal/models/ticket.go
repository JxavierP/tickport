package models

import (
	"time"

	"github.com/google/uuid"
)

type Priority string

const (
	PriorityLow    Priority = "Low"
	PriorityMedium Priority = "Medium"
	PriorityHigh   Priority = "High"
)

func (p Priority) IsValid() bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

type Status string

const (
	StatusOpen       Status = "Open"
	StatusInProgress Status = "In Progress"
	StatusClosed      Status = "Closed"
)

func (s Status) IsValid() bool {
	return s == StatusOpen || s == StatusInProgress || s == StatusClosed
}

type Ticket struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
