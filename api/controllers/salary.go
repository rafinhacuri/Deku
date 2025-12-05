package controllers

import (
	"sort"
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

type SalaryResponse struct {
	ID    int     `json:"id"`
	Vl    float64 `json:"vl" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Day   *int    `json:"day"`
	Month string  `json:"month" binding:"required"`
}

func GetSalaries(c *gin.Context) {
	userID := c.GetString("username")
	month := c.Query("month")

	rows, err := sqlite.SQL.Query("SELECT value, type, day, month FROM salary WHERE user = ? AND month = ?", userID, month)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	salaries := []SalaryResponse{}
	for rows.Next() {
		var salary SalaryResponse
		if err := rows.Scan(&salary.Vl, &salary.Type, &salary.Day, &salary.Month); err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
		salaries = append(salaries, salary)
	}

	sort.SliceStable(salaries, func(i, j int) bool {
		dayI := 0
		dayJ := 0
		if salaries[i].Day != nil {
			dayI = *salaries[i].Day
		}
		if salaries[j].Day != nil {
			dayJ = *salaries[j].Day
		}
		return dayI < dayJ
	})

	c.JSON(200, salaries)
}
