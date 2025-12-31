package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"sync_golang/config"
	"sync_golang/models"
	"sync_golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ================= REGISTER =================
func Register(c *gin.Context) {
	var input struct {
		Name     string
		Email    string
		Password string
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := models.User{
		Name:              input.Name,
		Email:             input.Email,
		Password:          utils.HashPassword(input.Password),
		Role:              "user",
		PasswordChangedAt: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register success"})
}

// ================= LOGIN =================
func Login(c *gin.Context) {
	var input struct {
		Email    string
		Password string
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak terdaftar"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan server"})
		}
		return
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Buat refresh token
	b := make([]byte, 32)
	rand.Read(b)
	refresh := hex.EncodeToString(b)

	config.DB.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token":  utils.CreateAccessToken(user),
		"refresh_token": refresh,
	})
}

// ================= REFRESH TOKEN =================
func Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var rt models.RefreshToken
	if err := config.DB.Where("token = ?", input.RefreshToken).First(&rt).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(rt.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, rt.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": utils.CreateAccessToken(user),
	})
}

// ================= CHANGE PASSWORD =================
func ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userCtx := c.MustGet("user").(models.User)

	var user models.User
	if err := config.DB.First(&user, userCtx.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if !utils.CheckPassword(user.Password, input.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah"})
		return
	}

	// Update password & reset refresh token
	config.DB.Model(&user).Updates(models.User{
		Password:          utils.HashPassword(input.NewPassword),
		PasswordChangedAt: time.Now(),
	})
	config.DB.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{})

	c.JSON(http.StatusOK, gin.H{"message": "Password diganti, silakan login ulang"})
}

// ================= LOGOUT =================
func Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	res := config.DB.Where("token = ?", input.RefreshToken).Delete(&models.RefreshToken{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}

// ================= ME =================
func Me(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
