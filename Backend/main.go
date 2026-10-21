package main

import (
	"jobtrack-backend/config"
	"jobtrack-backend/controllers"
	"jobtrack-backend/middleware"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()

	if config.DB != nil {
		config.DB.AutoMigrate(&models.User{}, &models.Job{})
	}

	router := gin.Default()

	// Auth Routes
	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)
	router.GET("/api/auth/profile", middleware.AuthMiddleware(), controllers.Profile)
	router.DELETE("/api/auth/delete-account", middleware.AuthMiddleware(), controllers.DeleteAccount)

	// Job CRUD Routes (Fitur 2)
	router.GET("/api/jobs", controllers.GetJobs)
	router.POST("/api/jobs", controllers.CreateJob)
	router.PUT("/api/jobs/:id", controllers.UpdateJob)
	router.DELETE("/api/jobs/:id", controllers.DeleteJob)

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
