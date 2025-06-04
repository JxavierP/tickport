package main

import (
	"time"

	"github.com/JxavierP/tickport/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://hoppscotch.io"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Api is Running",
		})
	})
	router.GET("/health", handlers.HealthCheckHandler)
	router.POST("/tickets", handlers.CreateTicketHandler)
	router.Run()
}
