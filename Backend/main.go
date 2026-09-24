package main

import (
	"jobtrack-backend/config"
	"jobtrack-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	// Memanggil semua rute API
	routes.SetupRoutes(r)

	r.Run(":8080")
}