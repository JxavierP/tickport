package main

import (
	"log"
	"time"

	"github.com/JxavierP/tickport/internal/database"
	"github.com/JxavierP/tickport/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB("tickport.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed %v", err)
	}
	defer db.Close()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://hoppscotch.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Api is Running",
		})
	})
	
	router.GET("/health", handlers.HealthCheckHandler(db))

	router.GET("/tickets", handlers.GetAllTicketsHandler(db))
	router.POST("/ticket", handlers.CreateTicketHandler(db))
	router.GET("/ticket/:id", handlers.GetTicketByIDHandler(db))
	router.PUT("/ticket/:id", handlers.UpdateTicketHandler(db))
	router.DELETE("/ticket/:id", handlers.DeleteTicketHandler(db))

	router.GET("/users", handlers.GetAllUsersHandler(db))
	router.GET("/user/:id", handlers.GetUserByIDHandler(db))
	router.POST("/user", handlers.CreateUserHandler(db))
	router.Run()
}
