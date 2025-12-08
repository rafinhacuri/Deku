package controllers

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/sqlite"
)

type Expense struct {
	Vl            float64 `json:"vl" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Description   *string `json:"description"`
	Day           *int    `json:"day"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
	Month         string  `json:"month" binding:"required"`
}

func InsertExpense(c *gin.Context) {
	var expense Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetString("username")

	_, err := sqlite.SQL.Exec("INSERT INTO expenses (user, value, type, day, month, paymentMethod, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID,
		expense.Vl,
		expense.Type,
		expense.Day,
		expense.Month,
		expense.PaymentMethod,
		expense.Description,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Expense inserted successfully"})
}

type ExpenseResponse struct {
	ID            int     `json:"id"`
	Vl            float64 `json:"vl" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Description   *string `json:"description"`
	Day           *int    `json:"day"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
	Month         string  `json:"month" binding:"required"`
}

func GetExpenses(c *gin.Context) {
	userID := c.GetString("username")
	month := c.Query("month")

	rows, err := sqlite.SQL.Query("SELECT id, value, type, day, paymentMethod, description, month FROM expenses WHERE user = ? AND month = ?", userID, month)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	var expenses []ExpenseResponse
	for rows.Next() {
		var exp ExpenseResponse
		if err := rows.Scan(&exp.ID, &exp.Vl, &exp.Type, &exp.Day, &exp.PaymentMethod, &exp.Description, &exp.Month); err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
		expenses = append(expenses, exp)
	}

	sort.SliceStable(expenses, func(i, j int) bool {
		var dayI, dayJ int
		if expenses[i].Day != nil {
			dayI = *expenses[i].Day
		} else {
			dayI = 32
		}
		if expenses[j].Day != nil {
			dayJ = *expenses[j].Day
		} else {
			dayJ = 32
		}
		return dayI < dayJ
	})

	for i, j := 0, len(expenses)-1; i < j; i, j = i+1, j-1 {
		expenses[i], expenses[j] = expenses[j], expenses[i]
	}

	c.JSON(200, expenses)
}

func DeleteExpense(c *gin.Context) {
	expenseID := c.Query("id")
	userID := c.GetString("username")

	result, err := sqlite.SQL.Exec("DELETE FROM expenses WHERE id = ? AND user = ?", expenseID, userID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Expense not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Expense deleted successfully"})
}

func UpdateExpense(c *gin.Context) {
	var expense Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	expenseID := c.Query("id")
	userID := c.GetString("username")

	result, err := sqlite.SQL.Exec("UPDATE expenses SET value = ?, type = ?, day = ?, paymentMethod = ?, description = ?, month = ?, updated_at = ? WHERE id = ? AND user = ?",
		expense.Vl,
		expense.Type,
		expense.Day,
		expense.PaymentMethod,
		expense.Description,
		expense.Month,
		time.Now().Format(time.RFC3339),
		expenseID,
		userID,
	)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Expense not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Expense updated successfully"})
}
