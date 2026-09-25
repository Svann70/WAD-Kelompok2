package main

import (
	"jobtrack-backend/src/config"
	"jobtrack-backend/src/middleware"

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

	//mulai server dengan port 3000
	router.Run(":3000")
}