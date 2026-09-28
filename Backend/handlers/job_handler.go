package handlers

import (
	"net/http"
	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
)

// Menampilkan semua job
func GetJobs(c *gin.Context) {
	var jobs []models.Job
	config.DB.Find(&jobs)
	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

// Menambahkan job baru
func CreateJob(c *gin.Context) {
	var input models.Job
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// Mengupdate job berdasarkan ID
func UpdateJob(c *gin.Context) {
	var job models.Job
	if err := config.DB.Where("id = ?", c.Param("id")).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan!"})
		return
	}

	var input models.Job
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&job).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": job})
}

// Menghapus job berdasarkan ID
func DeleteJob(c *gin.Context) {
	var job models.Job
	if err := config.DB.Where("id = ?", c.Param("id")).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan!"})
		return
	}

	config.DB.Delete(&job)
	c.JSON(http.StatusOK, gin.H{"data": true})
}