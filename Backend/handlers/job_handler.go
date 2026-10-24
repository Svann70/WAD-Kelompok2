package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// KONSTANTA & TIPE DATA
// ============================================================================

// ValidStatuses adalah satu-satunya sumber daftar status lamaran yang diakui
// aplikasi. Semua endpoint (create, update, update status) memvalidasi
// nilainya di sini, sehingga menambah status baru cukup dilakukan sekali.
var ValidStatuses = map[string]bool{
	"Terkirim":  true,
	"Interview": true,
	"Diterima":  true,
	"Ditolak":   true,
}

// pesanStatusInvalid adalah pesan error seragam untuk status di luar daftar.
const pesanStatusInvalid = "Status harus salah satu dari: Terkirim, Interview, Diterima, Ditolak"

// JobInput adalah bentuk payload JSON untuk create & update job.
// Field berbahasa Indonesia (perusahaan, posisi, dst.) disediakan sebagai
// alias agar payload gaya lama tetap diterima tanpa merusak yang baru.
type JobInput struct {
	// Field standar (dipakai frontend saat ini)
	Company     string `json:"company"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`

	// Alias bahasa Indonesia (kompatibilitas payload lama)
	Perusahaan string `json:"perusahaan"`
	Posisi     string `json:"posisi"`
	Tipe       string `json:"tipe"`
	Catatan    string `json:"catatan"`
}

// Normalize merapikan input: membuang spasi berlebih dan memetakan alias
// bahasa Indonesia ke field standar apabila field standarnya kosong.
func (in *JobInput) Normalize() {
	in.Company = firstNonEmpty(in.Company, in.Perusahaan)
	in.Title = firstNonEmpty(in.Title, in.Posisi)
	in.Type = firstNonEmpty(in.Type, in.Tipe)
	in.Description = firstNonEmpty(in.Description, in.Catatan)
	in.Link = strings.TrimSpace(in.Link)
	in.Status = strings.TrimSpace(in.Status)
}

// Validate memeriksa aturan bisnis untuk CREATE: field wajib terisi dan
// status (bila diberikan) harus valid. Mengembalikan pesan error,
// string kosong berarti input valid.
func (in *JobInput) Validate() string {
	if in.Company == "" || in.Title == "" {
		return "Field 'company'/'perusahaan' dan 'title'/'posisi' wajib diisi"
	}
	return in.ValidateStatus()
}

// ValidateStatus memeriksa status saja — dipakai oleh UPDATE yang sifatnya
// partial (field lain boleh tidak dikirim).
func (in *JobInput) ValidateStatus() string {
	if in.Status != "" && !ValidStatuses[in.Status] {
		return pesanStatusInvalid
	}
	return ""
}

// UpdateStatusInput adalah payload khusus update cepat status & link.
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
	Link   string `json:"link"`
}

// ============================================================================
// HELPER INTERNAL
// ============================================================================

// firstNonEmpty mengembalikan string pertama yang tidak kosong setelah
// di-trim; dipakai untuk memetakan alias field ke field standar.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// parseID mengubah parameter :id URL menjadi uint dan langsung membalas
// 400 Bad Request bila formatnya bukan angka positif.
func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid, ID harus berupa angka positif",
		})
		return 0, false
	}
	return uint(id), true
}

// dbGuard memastikan koneksi database tersedia sebelum handler menyentuh
// query. Mengembalikan false jika request sudah dijawab dengan error 500.
func dbGuard(c *gin.Context) bool {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Koneksi database tidak tersedia",
		})
		return false
	}
	return true
}

// catatRiwayatStatus menyimpan satu baris riwayat setiap kali status sebuah
// job berubah, sehingga progres lamaran bisa ditelusuri kembali.
func catatRiwayatStatus(jobID uint, status string) {
	config.DB.Create(&models.StatusHistory{
		JobID:     jobID,
		Status:    status,
		ChangedAt: time.Now(),
	})
}

// ============================================================================
// HANDLER JOB
// ============================================================================

// GetJobs menampilkan daftar job dengan pagination, pencarian, dan filter
// status. GET /api/jobs?page=1&limit=10&search=xxx&status=Terkirim
func GetJobs(c *gin.Context) {
	if !dbGuard(c) {
		return
	}

	// Parameter pagination dengan nilai default yang aman
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Bangun query bertahap: filter hanya ditambahkan bila parameternya ada
	query := config.DB.Model(&models.Job{})
	if search := c.Query("search"); search != "" {
		keyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR company ILIKE ?", keyword, keyword)
	}
	if status := c.Query("status"); status != "" && status != "Semua" {
		query = query.Where("status = ?", status)
	}

	// Hitung total item dulu (tanpa offset/limit) untuk kebutuhan pagination
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung data lamaran",
		})
		return
	}

	// Ambil data sesuai halaman yang diminta
	var jobs []models.Job
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data lamaran",
		})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jobs,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	})
}

// CreateJob menambahkan job baru. POST /api/jobs
func CreateJob(c *gin.Context) {
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
	if pesan := input.Validate(); pesan != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": pesan})
		return
	}

	if !dbGuard(c) {
		return
	}

	job := models.Job{
		Title:       input.Title,
		Company:     input.Company,
		Link:        input.Link,
		Type:        input.Type,
		Description: input.Description,
		Status:      input.Status,
		UserID:      input.UserID,
	}

	// Bila request melewati AuthMiddleware, user_id dari JWT menimpa nilai
	// dari payload sehingga job otomatis terikat pada pemilik yang login.
	// (claims JWT berupa angka selalu terbaca sebagai float64)
	if uid, exists := c.Get("user_id"); exists {
		if floatUID, ok := uid.(float64); ok {
			job.UserID = uint(floatUID)
		}
	}

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menambahkan data lamaran ke database",
			"error":   err.Error(),
		})
		return
	}

	// Status awal dicatat sebagai titik awal riwayat lamaran
	if job.Status != "" {
		catatRiwayatStatus(job.ID, job.Status)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Berhasil menambahkan lamaran baru",
		"data":    job,
	})
}

// UpdateJob mengedit job berdasarkan ID. PUT /api/jobs/:id
// Sifatnya partial: hanya field yang dikirim dan tidak kosong yang diubah.
func UpdateJob(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if !dbGuard(c) {
		return
	}

	var job models.Job
	if err := config.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data lamaran dengan ID tersebut tidak ditemukan",
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
	if pesan := input.ValidateStatus(); pesan != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": pesan})
		return
	}

	statusLama := job.Status

	// Terapkan hanya field yang diisi agar data lama tidak terhapus
	if input.Company != "" {
		job.Company = input.Company
	}
	if input.Title != "" {
		job.Title = input.Title
	}
	if input.Type != "" {
		job.Type = input.Type
	}
	if input.Link != "" {
		job.Link = input.Link
	}
	if input.Description != "" {
		job.Description = input.Description
	}
	if input.UserID != 0 {
		job.UserID = input.UserID
	}
	if input.Status != "" {
		job.Status = input.Status
	}

	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui data lamaran",
			"error":   err.Error(),
		})
		return
	}

	// Perubahan status lewat PUT pun tetap tercatat di riwayat
	if job.Status != statusLama {
		catatRiwayatStatus(job.ID, job.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil memperbarui data lamaran",
		"data":    job,
	})
}

// DeleteJob menghapus job berdasarkan ID (soft delete via GORM).
// DELETE /api/jobs/:id
func DeleteJob(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if !dbGuard(c) {
		return
	}

	var job models.Job
	if err := config.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data lamaran dengan ID tersebut tidak ditemukan",
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
		"data":    gin.H{"id": id},
	})
}

// UpdateJobStatus memperbarui status (dan opsional link) secara instan dari
// dropdown tabel. PATCH /api/jobs/:id/status
func UpdateJobStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if !dbGuard(c) {
		return
	}

	var job models.Job
	if err := config.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data job tidak ditemukan",
		})
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format payload status tidak valid",
		})
		return
	}

	if !ValidStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": pesanStatusInvalid})
		return
	}

	statusLama := job.Status
	updateData := map[string]interface{}{"status": input.Status}
	if input.Link != "" {
		updateData["link"] = input.Link
	}

	if err := config.DB.Model(&job).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui status",
		})
		return
	}

	// Selaraskan struct dengan hasil update untuk response
	job.Status = input.Status
	if input.Link != "" {
		job.Link = input.Link
	}

	// Riwayat hanya dicatat bila status benar-benar berubah
	if statusLama != input.Status {
		catatRiwayatStatus(job.ID, input.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status berhasil diperbarui",
		"data":    job,
	})
}

// GetJobStatusHistory menampilkan riwayat perubahan status sebuah job,
// diurutkan dari yang terbaru. GET /api/jobs/:id/status-history
func GetJobStatusHistory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if !dbGuard(c) {
		return
	}

	var histories []models.StatusHistory
	if err := config.DB.Where("job_id = ?", id).Order("changed_at desc").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": histories})
}
