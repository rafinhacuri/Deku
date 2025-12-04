package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/passwords"
	"github.com/rafinhacuri/Deku/sqlite"
)

func InsertUser(c *gin.Context) {
	var u user
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Invalid request payload"})
		return
	}

	if err := u.validate(); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	cryptPassword, cryptErr := passwords.BCrypt(u.Password)
	if cryptErr != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Failed to hash password"})
		return
	}

	insertUser := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Email:     u.Email,
		Password:  cryptPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	_, err := sqlite.SQL.Exec(
		"INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)",
		insertUser.Email,
		insertUser.Password,
		insertUser.CreatedAt,
		insertUser.UpdatedAt,
	)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}
