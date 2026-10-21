package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JobInput merepresentasikan format payload JSON untuk create & update job
type JobInput struct {
	Perusahaan  string `json:"perusahaan"`
	Posisi      string `json:"posisi"`
	Tipe        string `json:"tipe"`
	Tanggal     string `json:"tanggal"`
	Catatan     string `json:"catatan"`
	UserID      uint   `json:"user_id"`

	// Dukungan alias key (opsional) agar fleksibel jika frontend mengirim nama field alternatif
	Company     string `json:"company"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Normalize memetakan field alternatif ke field standar jika kosong
func (input *JobInput) Normalize() {
	if strings.TrimSpace(input.Perusahaan) == "" && strings.TrimSpace(input.Company) != "" {
		input.Perusahaan = strings.TrimSpace(input.Company)
	} else {
		input.Perusahaan = strings.TrimSpace(input.Perusahaan)
	}

	if strings.TrimSpace(input.Posisi) == "" && strings.TrimSpace(input.Title) != "" {
		input.Posisi = strings.TrimSpace(input.Title)
	} else {
		input.Posisi = strings.TrimSpace(input.Posisi)
	}

	if strings.TrimSpace(input.Tipe) == "" && strings.TrimSpace(input.Type) != "" {
		input.Tipe = strings.TrimSpace(input.Type)
	} else {
		input.Tipe = strings.TrimSpace(input.Tipe)
	}

	if strings.TrimSpace(input.Catatan) == "" && strings.TrimSpace(input.Description) != "" {
		input.Catatan = strings.TrimSpace(input.Description)
	} else {
		input.Catatan = strings.TrimSpace(input.Catatan)
	}

	input.Tanggal = strings.TrimSpace(input.Tanggal)
	if input.Tanggal == "" {
		input.Tanggal = time.Now().Format("2006-01-02")
	}
}

// GET /api/jobs
// Mengambil daftar lamaran dari database
func GetJobs(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Koneksi database tidak tersedia",
		})
		return
	}

	var jobs []models.Job
	if err := config.DB.Order("id desc").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil daftar lamaran",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar lamaran",
		"data":    jobs,
	})
}

// POST /api/jobs
// Menambahkan lamaran baru ke database
func CreateJob(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Koneksi database tidak tersedia",
		})
		return
	}

	var input JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payload JSON tidak valid",
			"error":   err.Error(),
		})
		return
	}

	input.Normalize()

	// Validasi field wajib
	if input.Perusahaan == "" || input.Posisi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Field 'perusahaan' dan 'posisi' wajib diisi",
		})
		return
	}

	job := models.Job{
		Perusahaan: input.Perusahaan,
		Posisi:     input.Posisi,
		Tipe:       input.Tipe,
		Tanggal:    input.Tanggal,
		Catatan:    input.Catatan,
		UserID:     input.UserID,
	}

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menambahkan data lamaran ke database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Berhasil menambahkan lamaran baru",
		"data":    job,
	})
}

// PUT /api/jobs/:id
// Mengedit data lamaran berdasarkan ID
func UpdateJob(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Koneksi database tidak tersedia",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid, ID harus berupa angka positif",
		})
		return
	}

	var job models.Job
	if err := config.DB.First(&job, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Data lamaran dengan ID tersebut tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mencari data lamaran",
			"error":   err.Error(),
		})
		return
	}

	var input JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payload JSON tidak valid",
			"error":   err.Error(),
		})
		return
	}

	input.Normalize()

	// Update field jika disediakan
	if input.Perusahaan != "" {
		job.Perusahaan = input.Perusahaan
	}
	if input.Posisi != "" {
		job.Posisi = input.Posisi
	}
	if input.Tipe != "" {
		job.Tipe = input.Tipe
	}
	if input.Tanggal != "" {
		job.Tanggal = input.Tanggal
	}
	if input.Catatan != "" {
		job.Catatan = input.Catatan
	}
	if input.UserID != 0 {
		job.UserID = input.UserID
	}

	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui data lamaran",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil memperbarui data lamaran",
		"data":    job,
	})
}

// DELETE /api/jobs/:id
// Menghapus data lamaran berdasarkan ID
func DeleteJob(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Koneksi database tidak tersedia",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid, ID harus berupa angka positif",
		})
		return
	}

	var job models.Job
	if err := config.DB.First(&job, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Data lamaran dengan ID tersebut tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mencari data lamaran",
			"error":   err.Error(),
		})
		return
	}

	if err := config.DB.Delete(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus data lamaran",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil menghapus data lamaran",
		"data": gin.H{
			"id": uint(id),
		},
	})
}
