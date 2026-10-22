package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
)

// Menampilkan job dengan pagination dan filter
func GetJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Job{})

	if search != "" {
		query = query.Where("title ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" && status != "Semua" {
		query = query.Where("status = ?", status)
	}

	var totalItems int64
	query.Count(&totalItems)

	var jobs []models.Job
	query.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs)

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": jobs,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	})
}

// Menambahkan job baru
func CreateJob(c *gin.Context) {
	var input models.Job
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uid, exists := c.Get("user_id"); exists {
		if floatUID, ok := uid.(float64); ok {
			input.UserID = uint(floatUID)
		}
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

// Menghapus job berdasarkan ID (soft delete)
func DeleteJob(c *gin.Context) {
	var job models.Job
	if err := config.DB.Where("id = ?", c.Param("id")).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan!"})
		return
	}

	config.DB.Delete(&job)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Struct khusus update status & link
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
	Link   string `json:"link"`
}

// UpdateJobStatus: Khusus memperbarui status & link secara instan
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

	updateData := map[string]interface{}{"status": input.Status}
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

// GetJobStatusHistory: Mengambil riwayat status job
func GetJobStatusHistory(c *gin.Context) {
	var histories []models.StatusHistory

	if err := config.DB.Where("job_id = ?", c.Param("id")).Order("changed_at desc").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat status!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": histories})
}
