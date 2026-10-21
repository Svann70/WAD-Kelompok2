package controllers

import (
	"net/http"
	"os"

	"jobtrack-backend/config"
	"jobtrack-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data registrasi tidak valid", "error": err.Error()})
		return
	}

	var existingUser models.User
	if result := config.DB.Where("email = ?", request.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengenkripsi password"})
		return
	}

	if request.Role == "" {
		request.Role = "user"
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     request.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat akun"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data login tidak valid", "error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Email atau password salah"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		},
	})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak terautentikasi"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak terautentikasi"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghapus akun"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Akun berhasil dihapus"})
}

// GetAllUsers: Endpoint admin untuk melihat seluruh user (contoh RBAC)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	config.DB.Select("id, name, email, role, created_at, updated_at").Find(&users)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
