package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("SECRET_KEY_GANTI_NANTI")

// ================= JWT AUTH =================
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token required"})
			return
		}

		tokenString := strings.Replace(auth, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}

// ================= ROLE CHECK =================
func RoleOnly(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "Forbidden",
		})
	}
}
