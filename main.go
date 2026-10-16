package main

import (
	"fmt"
	"log"
	"os"

	"github.com/enehizy/go-cms/models"
	"github.com/enehizy/go-cms/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("fialed to load .env file")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connected to database %w", err)
		return
	}
	fmt.Println("connected to database")
	mydb, err := db.DB()
	if err != nil {
		log.Fatalln("failed to connected to database %w", err)
		return
	}
	defer mydb.Close()
	if err := db.AutoMigrate(&models.Post{}, &models.Media{}, &models.Page{}); err != nil {
		log.Fatalln("failed to migrate database %w", err)
		return
	}
	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	{
		routes.CreateRoutes(api)

	}

	fmt.Println("server running at port:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("server stopped cause: %w", err)
	}

}
