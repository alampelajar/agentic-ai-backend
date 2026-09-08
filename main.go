package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"agentic-ai-backend/config"
	"agentic-ai-backend/handlers"
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