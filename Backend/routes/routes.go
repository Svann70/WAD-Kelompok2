package routes

import (
	"jobtrack-backend/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/jobs", handlers.GetJobs)
		api.POST("/jobs", handlers.CreateJob)
		api.PUT("/jobs/:id", handlers.UpdateJob)
		api.DELETE("/jobs/:id", handlers.DeleteJob)
	}
}