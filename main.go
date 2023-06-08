package main

import (
	"sarkor_telekom/database"
	"sarkor_telekom/handlers"
	"sarkor_telekom/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create database
	database.InitDB()

	router := gin.Default()

	router.POST("/user/register", handlers.RegisterUser)
	router.POST("/user/auth", handlers.AuthenticateUser)

	// Register middleware
	router.Use(middleware.AuthMiddleware)

	router.GET("/user/:name", handlers.GetUserByName)
	router.POST("/user/phone", handlers.AddPhoneNumber)
	router.GET("/user/phone", handlers.GetUsersByPhoneNumber)
	router.PUT("/user/phone", handlers.UpdatePhoneNumber)
	router.DELETE("/user/phone/:phone_id", handlers.DeletePhoneNumber)

	// Start the server
	router.Run(":8080")
}
