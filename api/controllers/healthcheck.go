package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/sqlite"
)

func HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	if err := sqlite.TestSqlite(ctx); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok!",
	})
}
