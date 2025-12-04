package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafinhacuri/Deku/env"
	"github.com/rafinhacuri/Deku/routes"
	"github.com/rafinhacuri/Deku/sqlite"
)

func init() {
	godotenv.Load("../.env")

	env.LoadEnv()

	if _, err := sqlite.InitSqlite(); err != nil {
		panic(err)
	}

}

func main() {

	gin.DefaultWriter = io.Discard

	server := gin.Default()

	server.Use(gin.LoggerWithWriter(os.Stdout, "/healthcheck"))

	server.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	routes.RegisterRoutes(server)
	server.Run(":8080")
}
