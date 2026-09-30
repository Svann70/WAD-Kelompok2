package main

import (
	"jobtrack-backend/config"
	"jobtrack-backend/middleware"
	"jobtrack-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()

	router := gin.Default()

	router.Use(middleware.SetupCORS())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Server Jobtrack Berhasil Berjalan",	
		})
	})

	routes.SetupRoutes(router)
	
	//mulai server dengan port 3000
	router.Run(":3000")
}