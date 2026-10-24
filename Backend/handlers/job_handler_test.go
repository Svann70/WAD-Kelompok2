package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupTestRouter membuat router Gin ringan tanpa database. Berkat dbGuard
// di setiap handler, semua pengujian validasi input bisa berjalan tanpa
// koneksi PostgreSQL yang sebenarnya.
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/jobs", GetJobs)
	r.POST("/api/jobs", CreateJob)
	r.PUT("/api/jobs/:id", UpdateJob)
	r.DELETE("/api/jobs/:id", DeleteJob)
	r.PATCH("/api/jobs/:id/status", UpdateJobStatus)
	r.GET("/api/jobs/:id/status-history", GetJobStatusHistory)
	return r
}

// performRequest membantu mengirim request ke router test dan mengembalikan
// response recorder-nya.
func performRequest(r *gin.Engine, method, path string, payload []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if payload != nil {
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Payload tanpa company/perusahaan dan title/posisi harus ditolak 400
// sebelum menyentuh database.
func TestCreateJob_ValidationError(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodPost, "/api/jobs", []byte(`{"tipe":"Magang"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// Alias bahasa Indonesia harus diterima: payload ini lolos validasi,
// lalu berhenti di dbGuard (500) karena test berjalan tanpa database.
func TestCreateJob_AliasFieldIndonesia(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodPost, "/api/jobs",
		[]byte(`{"perusahaan":"PT ABC","posisi":"Backend Developer"}`))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 (lolos validasi, berhenti di dbGuard), got %d. Body: %s", w.Code, w.Body.String())
	}
}

// Status di luar daftar ValidStatuses harus ditolak 400.
func TestCreateJob_InvalidStatus(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodPost, "/api/jobs",
		[]byte(`{"company":"PT ABC","title":"Backend","status":"Sembarangan"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// ID bukan angka harus ditolak 400, bukan 404.
func TestUpdateJob_InvalidID(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodPut, "/api/jobs/invalid-id",
		[]byte(`{"company":"PT ABC","posisi":"Backend"}`))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// ID bernilai 0 juga tidak valid.
func TestDeleteJob_InvalidID(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodDelete, "/api/jobs/0", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// Endpoint riwayat status juga menolak ID non-numerik dengan 400.
func TestGetJobStatusHistory_InvalidID(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodGet, "/api/jobs/abc/status-history", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// Tanpa koneksi database, GET /api/jobs harus dijawab 500 yang rapi
// (bukan panic), sesuai kontrak dbGuard.
func TestGetJobs_DatabaseUnavailable(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, http.MethodGet, "/api/jobs", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d. Body: %s", w.Code, w.Body.String())
	}
}
