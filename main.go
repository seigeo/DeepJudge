package main

import (
	"deepjudge/routes"
	"deepjudge/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
