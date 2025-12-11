package controllers

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/sqlite"
)

type Income struct {
	Vl    float64 `json:"vl" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Day   *int    `json:"day"`
	Month string  `json:"month" binding:"required"`
}

func InsertIncome(c *gin.Context) {
	var income Income
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetString("username")

	_, err := sqlite.SQL.Exec("INSERT INTO income (user, value, type, day, month, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID,
		income.Vl,
		income.Type,
		income.Day,
		income.Month,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Income inserted successfully"})
}

type IncomeResponse struct {
	ID    int     `json:"id"`
	Vl    float64 `json:"vl" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Day   *int    `json:"day"`
	Month string  `json:"month" binding:"required"`
}

func GetIncomes(c *gin.Context) {
	userID := c.GetString("username")
	month := c.Query("month")

	rows, err := sqlite.SQL.Query("SELECT id, value, type, day, month FROM income WHERE user = ? AND month = ?", userID, month)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	incomes := []IncomeResponse{}
	for rows.Next() {
		var income IncomeResponse
		if err := rows.Scan(&income.ID, &income.Vl, &income.Type, &income.Day, &income.Month); err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
		incomes = append(incomes, income)
	}

	sort.SliceStable(incomes, func(i, j int) bool {
		dayI := 0
		dayJ := 0
		if incomes[i].Day != nil {
			dayI = *incomes[i].Day
		}
		if incomes[j].Day != nil {
			dayJ = *incomes[j].Day
		}
		return dayI < dayJ
	})

	for i, j := 0, len(incomes)-1; i < j; i, j = i+1, j-1 {
		incomes[i], incomes[j] = incomes[j], incomes[i]
	}

	c.JSON(200, incomes)
}

func DeleteIncome(c *gin.Context) {
	id := c.Query("id")

	res, err := sqlite.SQL.Exec("DELETE FROM income WHERE id = ?", id)
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

func UpdateIncome(c *gin.Context) {
	id := c.Query("id")
	var income Income
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	row, err := sqlite.SQL.Exec("UPDATE income SET value = ?, type = ?, day = ?, updated_at = ? WHERE id = ?",
		income.Vl,
		income.Type,
		income.Day,
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
		c.JSON(404, gin.H{"message": "Income not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Income updated successfully"})
}
