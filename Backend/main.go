package main

import (
	"jobtrack-backend/config"
	"jobtrack-backend/routes"
	"jobtrack-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	// Memanggil semua rute API
	routes.SetupRoutes(r)

	r.Run(":8080")

	config.ConnectDatabase()

	router := gin.Default()

	router.Use(middleware.SetupCORS())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Server Jobtrack Berhasil Berjalan",	
		})
	})

	//mulai server dengan port 3000
	router.Run(":3000")
}