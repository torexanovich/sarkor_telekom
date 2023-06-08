package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("SESSTOKEN")
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized."})
		c.Abort()
		return
	}

	fmt.Println("COOKIE>", cookie)

	c.Next()
}
