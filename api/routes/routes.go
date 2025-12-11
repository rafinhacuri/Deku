package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/Deku/controllers"
	"github.com/rafinhacuri/Deku/middleware"
)

func RegisterRoutes(server *gin.Engine) {
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "Not Found",
		})
	})

	server.GET("/healthcheck", controllers.HealthCheck)
	server.PUT("/api/user", controllers.InsertUser)

	server.POST("/login", controllers.Auth)

	api := server.Group("/api", middleware.Authenticate)
	api.GET("/check-session", middleware.CheckSession)

	api.PUT("/income", controllers.InsertIncome)
	api.GET("/income", controllers.GetIncomes)
	api.DELETE("/income", controllers.DeleteIncome)
	api.POST("/income", controllers.UpdateIncome)

	api.PUT("/expense", controllers.InsertExpense)
	api.GET("/expense", controllers.GetExpenses)
	api.DELETE("/expense", controllers.DeleteExpense)
	api.POST("/expense", controllers.UpdateExpense)
}
