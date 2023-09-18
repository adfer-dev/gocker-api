package main

import (
	"gocker-api/routes"

	"github.com/gin-gonic/gin"
)

func initRoutes(app *gin.Engine) {
	routes.InitUserRoutes(app)
	routes.InitAuthRoutes(app)
}

func main() {
	gin.SetMode("debug")

	app := gin.Default()
	app.ForwardedByClientIP = true
	app.SetTrustedProxies([]string{"127.0.0.1"})

	initRoutes(app)
	app.Run(":8080")
}
