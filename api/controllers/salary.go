package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/sqlite"
)

type Salary struct {
	Vl    float64 `json:"vl" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Day   *int    `json:"day"`
	Month string  `json:"month" binding:"required"`
}

func InsertSalary(c *gin.Context) {
	var salary Salary
	if err := c.ShouldBindJSON(&salary); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetString("username")

	_, err := sqlite.SQL.Exec("INSERT INTO salary (user, value, type, day, month, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID,
		salary.Vl,
		salary.Type,
		salary.Day,
		salary.Month,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Salary inserted successfully"})
}
