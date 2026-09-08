package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"agentic-ai-backend/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke MySQL:", err)
	}

	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migration:", err)
	}

	fmt.Println("MySQL berhasil terhubung")
	fmt.Println("Database migration berhasil")
}