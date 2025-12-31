package middleware

import (
	"net/http"
	"strings"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			return
		}

		tokenString := strings.Replace(auth, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		user := models.User{
			ID:   uint(claims["user_id"].(float64)),
			Role: claims["role"].(string),
		}

		if len(roles) > 0 {
			allowed := false
			for _, r := range roles {
				if r == user.Role {
					allowed = true
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}
		}

		c.Set("user", user)
		c.Next()
	}
}
