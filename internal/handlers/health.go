package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := db.Ping(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "Database unreachable"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}
