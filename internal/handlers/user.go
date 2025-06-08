package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/JxavierP/tickport/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string      `json:"username" binding:"required"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=8"`
	Role     models.Role `json:"role" binding:"required,oneof=admin agent"`
}

func GetAllUsersHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.Query("role")
		if role != "" && role != string(models.AdminRole) && role != string(models.AgentRole) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}

		query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE 1=1`
		args := []interface{}{}
		if role != "" {
			query += ` AND role = ?`
			args = append(args, role)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			var idStr string
			var roleStr models.Role
			if err := rows.Scan(
				&idStr,
				&user.Username,
				&user.Email,
				&roleStr,
				&user.CreatedAt,
				&user.UpdatedAt,
			); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user data"})
				return
			}

			user.ID, _ = uuid.Parse(idStr)
			if roleStr.IsValid() {
				user.Role = models.Role(roleStr)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role in database"})
				return
			}

			users = append(users, user)
		}

		ctx.JSON(http.StatusOK, users)
	}
}

func GetUserByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = ?`
		var user models.User
		var idStr string
		var roleStr models.Role

		if err := db.QueryRow(query, id).Scan(
			&idStr,
			&user.Username,
			&user.Email,
			&roleStr,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		user.ID, _ = uuid.Parse(idStr)
		if roleStr.IsValid() {
			user.Role = models.Role(roleStr)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role in database"})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func CreateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		user := models.User{
			ID:       uuid.New(),
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password, // Password should be hashed before storing in production
			Role:     req.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		query := `INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
		if _, err := db.Exec(query,
			user.ID,
			user.Username,
			user.Email,
			user.Password,
			user.Role,
			user.CreatedAt,
			user.UpdatedAt,
		); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
	}
}

// Todo: Implement UpdateUserHandler and DeleteUserHandler
