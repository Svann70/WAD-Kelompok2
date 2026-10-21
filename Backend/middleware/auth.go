package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware: Memvalidasi JWT dan menyimpan user_id & role ke context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Format token tidak valid"})
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token tidak valid atau sudah kadaluarsa"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Claims token tidak valid"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		if role, exists := claims["role"]; exists {
			c.Set("role", role)
		}

		c.Next()
	}
}

// RoleMiddleware: Memeriksa role pengguna (RBAC)
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Role tidak ditemukan"})
			c.Abort()
			return
		}

		userRole, _ := roleVal.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Akses ditolak: role tidak memiliki izin"})
		c.Abort()
	}
}
