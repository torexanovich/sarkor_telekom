package main

import (
	"sarkor_telekom/database"
	"sarkor_telekom/handlers"
	"sarkor_telekom/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	router := gin.Default()

	router.POST("/user/register", handlers.RegisterUser)
	router.POST("/user/auth", handlers.AuthenticateUser)

	router.Use(middleware.AuthMiddleware)

	router.GET("/user/:name", handlers.GetUserByName)
	router.POST("/user/phone", handlers.AddPhoneNumber)
	router.GET("/user/phone", handlers.GetUsersByPhoneNumber)
	router.PUT("/user/phone", handlers.UpdatePhoneNumber)
	router.DELETE("/user/phone/:phone_id", handlers.DeletePhoneNumber)

	router.Run(":8080")	
}
