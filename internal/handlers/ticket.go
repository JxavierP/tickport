package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/JxavierP/tickport/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
}

func CreateTicketHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateTicketRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}

		if req.Title == "" || req.Description == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title and Description are required"})
			return
		}

		priority := models.Priority(req.Priority)
		status := models.Status(req.Status)

		if !priority.IsValid() || !status.IsValid() {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority or status"})
			return
		}

		ticket := models.Ticket{
			ID:          uuid.New(),
			Title:       req.Title,
			Description: req.Description,
			Priority:    priority,
			Status:      status,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		query := `
			INSERT INTO tickets (id, title, description, priority, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`

		_, err := db.Exec(query,
			ticket.ID.String(),
			ticket.Title,
			ticket.Description,
			ticket.Priority,
			ticket.Status,
			ticket.CreatedAt,
			ticket.UpdatedAt,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ticket"})
		}
		
		ctx.JSON(http.StatusCreated, ticket)
	}
}
