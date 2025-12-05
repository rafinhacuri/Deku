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

	rows, err := sqlite.SQL.Query("SELECT id, value, type, day, month FROM salary WHERE user = ? AND month = ?", userID, month)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	salaries := []SalaryResponse{}
	for rows.Next() {
		var salary SalaryResponse
		if err := rows.Scan(&salary.ID, &salary.Vl, &salary.Type, &salary.Day, &salary.Month); err != nil {
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

	for i, j := 0, len(salaries)-1; i < j; i, j = i+1, j-1 {
		salaries[i], salaries[j] = salaries[j], salaries[i]
	}

	c.JSON(200, salaries)
}

func DeleteSalary(c *gin.Context) {
	id := c.Query("id")

	res, err := sqlite.SQL.Exec("DELETE FROM salary WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Income not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Income deleted successfully"})
}

func UpdateSalary(c *gin.Context) {
	id := c.Query("id")
	var salary Salary
	if err := c.ShouldBindJSON(&salary); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	row, err := sqlite.SQL.Exec("UPDATE salary SET value = ?, type = ?, day = ?, updated_at = ? WHERE id = ?",
		salary.Vl,
		salary.Type,
		salary.Day,
		time.Now().Format(time.RFC3339),
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Salary not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Salary updated successfully"})
}
