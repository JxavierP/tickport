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

func GetAllTicketsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := db.Query(`SELECT * FROM tickets`)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		}
		defer rows.Close()

		var tickets []models.Ticket

		for rows.Next() {
			var t models.Ticket
			var idStr string
			if err := rows.Scan(
				&idStr,
				&t.Title,
				&t.Description,
				&t.Priority,
				&t.Status,
				&t.CreatedAt,
				&t.UpdatedAt,
			); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan ticket"})
				return
			}

			t.ID, err = uuid.Parse(idStr)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID"})
				return
			}

			tickets = append(tickets, t)
		}
		ctx.JSON(http.StatusOK, tickets)
	}
}

func GetTicketByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		var t models.Ticket
		query := `
            SELECT id, title, description, priority, status, created_at, updated_at
            FROM tickets
            WHERE id = ?
        `
		row := db.QueryRow(query, id.String())

		var idRaw string
		err = row.Scan(
			&idRaw,
			&t.Title,
			&t.Description,
			&t.Priority,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket"})
			return
		}

		t.ID, _ = uuid.Parse(idRaw)

		c.JSON(http.StatusOK, t)
	}
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

func UpdateTicketHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

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
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			Priority:    priority,
			Status:      status,
			UpdatedAt:   time.Now(),
		}

		query := `
			UPDATE tickets
			SET title = ?, description = ?, priority = ?, status = ?, updated_at = ?
			WHERE id = ?
		`

		result, err := db.Exec(query,
			ticket.Title,
			ticket.Description,
			ticket.Priority,
			ticket.Status,
			ticket.UpdatedAt,
			ticket.ID.String(),
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		ctx.JSON(http.StatusOK, ticket)
	}
}

func DeleteTicketHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		query := `DELETE FROM tickets WHERE id = ?`
		result, err := db.Exec(query, id.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
			return
		}

		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
	}
}
