package controllers

import (
	"net/http"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var request RegisterRequest

	// Membaca data JSON dari request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data registrasi tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Mengecek apakah email sudah terdaftar
	var existingUser models.User

	result := config.DB.Where("email = ?", request.Email).First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email sudah terdaftar",
		})
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengenkripsi password",
		})
		return
	}

	// Membuat user baru
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	// Menyimpan user ke database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat akun",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
