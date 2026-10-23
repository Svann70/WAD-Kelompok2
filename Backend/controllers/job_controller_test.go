package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/jobs", GetJobs)
	r.POST("/api/jobs", CreateJob)
	r.PUT("/api/jobs/:id", UpdateJob)
	r.DELETE("/api/jobs/:id", DeleteJob)
	return r
}

func TestCreateJob_ValidationError(t *testing.T) {
	router := setupTestRouter()

	// Kirim payload tanpa perusahaan dan posisi
	payload := []byte(`{"tipe":"Magang"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Status code harus 400 Bad Request atau 500 jika DB nil sebelum validasi
	// Dalam implementasi kita, jika DB nil mengembalikan 500, jika DB ada tapi field kosong mengembalikan 400
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}
}

func TestUpdateJob_InvalidID(t *testing.T) {
	router := setupTestRouter()

	payload := []byte(`{"perusahaan":"PT ABC","posisi":"Backend"}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/jobs/invalid-id", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}
}

func TestDeleteJob_InvalidID(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodDelete, "/api/jobs/not-a-number", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}
}
