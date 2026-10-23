package main

import (
	"os"

	"jobtrack-backend/config"
	"jobtrack-backend/controllers"
	"jobtrack-backend/middleware"
	"jobtrack-backend/models"
	"jobtrack-backend/routes"
	"jobtrack-backend/seeders"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	if config.DB != nil {
		config.DB.AutoMigrate(&models.User{}, &models.Job{}, &models.StatusHistory{})
	}

	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		seeders.SeedAll()
		return
	}

	router := gin.Default()
	router.Use(middleware.SetupCORS())

	// Endpoint Auth
	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)
	router.GET("/api/auth/profile", middleware.AuthMiddleware(), controllers.Profile)
	router.DELETE("/api/auth/delete-account", middleware.AuthMiddleware(), controllers.DeleteAccount)

	// Endpoint Admin (RBAC: khusus role admin)
	router.GET("/api/admin/users", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetAllUsers)

	// Healthcheck
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Server Jobtrack Berhasil Berjalan",
		})
	})

	// Setup routes (Job CRUD, Quick Status, Status History)
	routes.SetupRoutes(router)

	router.Run(":3000")
}
