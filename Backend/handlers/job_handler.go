package handlers

import (
	"net/http"
	"time"
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

// Struct khusus untuk membaca input status & link
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
	Link   string `json:"link"`
}

// UpdateJobStatus: Khusus memperbarui status (Terkirim, Interview, Diterima, Ditolak) & Link secara instan
func UpdateJobStatus(c *gin.Context) {
	var job models.Job

	if err := config.DB.Where("id = ?", c.Param("id")).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data job tidak ditemukan!"})
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status tidak valid!"})
		return
	}

	validStatuses := map[string]bool{
		"Terkirim":  true,
		"Interview": true,
		"Diterima":  true,
		"Ditolak":   true,
	}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus salah satu dari: Terkirim, Interview, Diterima, Ditolak"})
		return
	}

	updateData := map[string]interface{}{
		"status": input.Status,
	}
	if input.Link != "" {
		updateData["link"] = input.Link
	}

	config.DB.Model(&job).Updates(updateData)

	history := models.StatusHistory{
		JobID:     job.ID,
		Status:    input.Status,
		ChangedAt: time.Now(),
	}
	config.DB.Create(&history)

	c.JSON(http.StatusOK, gin.H{
		"message": "Status berhasil diperbarui!",
		"data":    job,
	})
}

func GetJobStatusHistory(c *gin.Context) {
	var histories []models.StatusHistory

	if err := config.DB.Where("job_id = ?", c.Param("id")).Order("changed_at desc").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat status!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": histories})
}