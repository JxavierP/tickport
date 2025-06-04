package handlers

import (
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

func CreateTicketHandler(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if req.Title == "" || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Description are required"})
		return
	}

	priority := models.Priority(req.Priority)
	status := models.Status(req.Status)

	if !priority.IsValid() || !status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority or status"})
		return
	}

	ticketID := uuid.New()
	ticket := models.Ticket{
		ID:          ticketID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusCreated, ticket)
}
