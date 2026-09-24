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

// Struct khusus untuk membaca input status & link
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
	Link   string `json:"link"`
}

// UpdateJobStatus: Khusus memperbarui status (Terkirim, Interview, Diterima, Ditolak) & Link secara instan
func UpdateJobStatus(c *gin.Context) {
	var job models.Job

	// 1. Cari job berdasarkan ID
	if err := config.DB.Where("id = ?", c.Param("id")).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data job tidak ditemukan!"})
		return
	}

	// 2. Validasi input JSON dari Request Body
	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status tidak valid!"})
		return
	}

	// 3. Buat map untuk menampung field yang di-update
	updateData := map[string]interface{}{
		"status": input.Status,
	}
	if input.Link != "" {
		updateData["link"] = input.Link
	}

	// 4. Update data di database (GORM otomatis memperbarui field UpdatedAt)
	config.DB.Model(&job).Updates(updateData)

	// 5. Kirim respon balik ke Client/Frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Status berhasil diperbarui!",
		"data":    job,
	})
}