package utils

import (
	"time"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/golang-jwt/jwt"
)

func CreateAccessToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	t, _ := token.SignedString(config.JwtKey)
	return t
}
