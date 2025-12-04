package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/utils"
)

func Authenticate(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	email, err := utils.JWTValidate(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	c.Set("username", email)
	c.Next()
}

func CheckSession(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	email, err := utils.JWTValidate(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	c.JSON(200, gin.H{"message": email})
}
