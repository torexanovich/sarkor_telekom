package handlers

import (
	"fmt"
	"sarkor_telekom/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	name := c.PostForm("name")
	ages := c.PostForm("age")
	age, _ := strconv.Atoi(ages)
	err := database.CreateUser(login, password, name, age)
	if err != nil {
		fmt.Println("Error creating user:", err)
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func AuthenticateUser(c *gin.Context) {
	var claims struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err := c.ShouldBindJSON(&claims)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	token, userID, err := database.AuthenticateUser(claims.Login, claims.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	c.SetCookie("SESSTOKEN", token, 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "Authentication successful",
		"token":   token,
		"userID":  userID,
	})
}

func GetUserByName(c *gin.Context) {
	name := c.Param("name") // /user/John

	user, err := database.GetUser(name)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(200, gin.H{"id": user.ID, "name": user.Name, "age": user.Age})
}

func AddPhoneNumber(c *gin.Context) {
	var Phone struct {
		PhoneNumber string `json:"phone"`
		Description string `json:"description"`
		IsMobile    bool   `json:"is_mobile"`
	}
	err := c.ShouldBindJSON(&Phone)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := database.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	err = database.CreatePhoneNumber(Phone.PhoneNumber, Phone.Description, Phone.IsMobile, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add phone number"})
		return
	}

	c.JSON(200, gin.H{"message": "Phone number added successfully"})
}

func GetUsersByPhoneNumber(c *gin.Context) {
	q := c.Query("q") // /user/phone?q=1234

	users, err := database.GetUsersByPhoneNumber(q)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get users by number"})
		return
	}

	c.JSON(200, users)
}

func UpdatePhoneNumber(c *gin.Context) {
	var Phone struct {
		PhoneID     int    `json:"phone_id"`
		PhoneNumber string `json:"phone"`
		Description string `json:"description"`
		IsMobile    bool   `json:"is_mobile"`
	}
	err := c.ShouldBindJSON(&Phone)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := database.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	err = database.UpdatePhoneNumber(Phone.PhoneID, Phone.PhoneNumber, Phone.Description, Phone.IsMobile, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update phone number"})
		return
	}

	c.JSON(200, gin.H{"message": "Phone number updated successfully"})
}

func DeletePhoneNumber(c *gin.Context) {
	phoneIDs := c.Param("phone_id")
	phoneID, _ := strconv.Atoi(phoneIDs)

	userID, err := database.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	err = database.DeletePhoneNumber(phoneID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete phone number"})
		return
	}

	c.JSON(200, gin.H{"message": "Phone number deleted successfully"})
}
