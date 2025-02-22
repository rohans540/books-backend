package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rohans540/books-backend/database"
	_ "github.com/rohans540/books-backend/docs"
	"github.com/rohans540/books-backend/kafka"
	"github.com/rohans540/books-backend/redis"
	"github.com/rohans540/books-backend/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Books API
// @version 1.0
// @description This is a simple API for managing books.
// @host localhost:8000
// @BasePath /

func main() {
	// Load environment variables
	database.ConnectDB()
	kafka.InitProducer()
	redis.ConnectRedis()

	router := gin.Default()

	router.Use(cors.Default())

	// Swagger Documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Routes
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(router.Run(":" + port))
}
