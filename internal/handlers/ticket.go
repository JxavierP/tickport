package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/JxavierP/tickport/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Priority    models.Priority `json:"priority" binding:"required"`
	Status      models.Status   `json:"status" binding:"required"`
	CreatorID   string          `json:"creator_id"`
	AssigneeID  *string         `json:"assignee_id"`
}

func GetAllTicketsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := ctx.Query("status")
		priority := ctx.Query("priority")
		sort := ctx.DefaultQuery("sort", "created_at")
		order := ctx.DefaultQuery("order", "desc")
		creatorID := ctx.Query("creator_id")
		assigneeID := ctx.Query("assignee_id")

		if status != "" && !models.Status(status).IsValid() {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		if priority != "" && !models.Priority(priority).IsValid() {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority"})
			return
		}

		allowedSorts := map[string]bool{
			"created_at": true, "updated_at": true, "title": true,
			"priority": true, "status": true, "creator_id": true, "assignee_id": true,
		}

		if !allowedSorts[sort] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort field"})
			return
		}

		if order != "asc" && order != "desc" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order"})
			return
		}

		query := `SELECT id, title, description, priority, status, creator_id, assignee_id, created_at, updated_at FROM tickets WHERE 1=1`
		args := []interface{}{}
		if status != "" {
			query += ` AND status = ?`
			args = append(args, status)
		}
		if priority != "" {
			query += ` AND priority = ?`
			args = append(args, priority)
		}
		if creatorID != "" {
			if _, err := uuid.Parse(creatorID); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creator ID"})
				return
			}
			query += ` AND creator_id = ?`
			args = append(args, creatorID)
		}
		if assigneeID != "" {
			if _, err := uuid.Parse(assigneeID); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
				return
			}
			query += ` AND assignee_id = ?`
			args = append(args, assigneeID)
		}
		query += fmt.Sprintf(` ORDER BY %s %s`, sort, order)

		rows, err := db.Query(query, args...)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
			return
		}
		defer rows.Close()

		var tickets []models.Ticket
		for rows.Next() {
			var t models.Ticket
			var idStr, creatorStr string
			var assigneeStr sql.NullString
			if err := rows.Scan(
				&idStr,
				&t.Title,
				&t.Description,
				&t.Priority,
				&t.Status,
				&creatorStr,
				&assigneeStr,
				&t.CreatedAt,
				&t.UpdatedAt); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan ticket"})
				return
			}

			t.ID, _ = uuid.Parse(idStr)
			t.CreatorID, _ = uuid.Parse(creatorStr)

			if assigneeStr.Valid {
				id, _ := uuid.Parse(assigneeStr.String)
				t.AssigneeID = &id
			} else {
				t.AssigneeID = nil
			}
			tickets = append(tickets, t)
		}

		ctx.JSON(http.StatusOK, tickets)
	}
}

func GetTicketByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		var t models.Ticket
		var assignee sql.NullString

		query := `SELECT id, title, description, priority, status, creator_id, assignee_id, created_at, updated_at 
		FROM tickets WHERE id = ?`
		row := db.QueryRow(query, id.String())

		var idRaw, creatorRaw string
		err = row.Scan(
			&idRaw,
			&t.Title,
			&t.Description,
			&t.Priority,
			&t.Status,
			&creatorRaw,
			&assignee,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket"})
			return
		}

		t.ID, _ = uuid.Parse(idRaw)
		t.CreatorID, _ = uuid.Parse(creatorRaw)

		if assignee.Valid {
			id, _ := uuid.Parse(assignee.String)
			t.AssigneeID = &id
		}

		ctx.JSON(http.StatusOK, t)
	}
}

func CreateTicketHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateTicketRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}

		creatorID, err := uuid.Parse(req.CreatorID)
		if err != nil || creatorID == uuid.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or empty creator ID"})
			return
		}

		var assignee sql.NullString
		var assigneeID *uuid.UUID
		if req.AssigneeID != nil {
			id, err := uuid.Parse(*req.AssigneeID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
				return
			}
			assignee = sql.NullString{String: id.String(), Valid: true}
			assigneeID = &id
		}

		ticket := models.Ticket{
			ID:          uuid.New(),
			Title:       req.Title,
			Description: req.Description,
			Priority:    req.Priority,
			Status:      req.Status,
			CreatorID:   creatorID,
			AssigneeID:  assigneeID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		query := `INSERT INTO tickets 
		(id, title, description, priority, status, creator_id, assignee_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err = db.Exec(query,
			ticket.ID,
			ticket.Title,
			ticket.Description,
			ticket.Priority,
			ticket.Status,
			ticket.CreatorID,
			assignee,
			ticket.CreatedAt,
			ticket.UpdatedAt,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ticket"})
			return
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

		var assigneeID sql.NullString
		if req.AssigneeID != nil {
			assigneeUUID, err := uuid.Parse(*req.AssigneeID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
				return
			}
			assigneeID = sql.NullString{String: assigneeUUID.String(), Valid: true}
		}

		query := `UPDATE tickets 
		SET title = ?, description = ?, priority = ?, status = ?, assignee_id = ?, updated_at = ? 
		WHERE id = ?`

		result, err := db.Exec(query, req.Title, req.Description, req.Priority, req.Status, assigneeID, time.Now(), id.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
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
