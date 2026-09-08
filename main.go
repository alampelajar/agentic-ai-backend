package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"agentic-ai-backend/config"
	"agentic-ai-backend/handlers"
	"agentic-ai-backend/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: file .env tidak ditemukan")
	}

	config.ConnectDatabase()

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.Refresh)
			auth.POST("/logout", handlers.Logout)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handlers.Me)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Backend berjalan di http://localhost:" + port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}