package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"agentic-ai-backend/config"
	"agentic-ai-backend/models"
	"agentic-ai-backend/utils"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Nama wajib diisi",
		})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Email wajib diisi",
		})
		return
	}

	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Password minimal 8 karakter",
		})
		return
	}

	var existingUser models.User

	result := config.DB.
		Where("email = ?", req.Email).
		First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"ok":      false,
			"message": "Email sudah terdaftar",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"message": "Gagal mengenkripsi password",
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"message": "Gagal membuat pengguna",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok":      true,
		"message": "Registrasi berhasil",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "Email dan password wajib diisi",
		})
		return
	}

	var user models.User

	result := config.DB.
		Where("email = ?", req.Email).
		First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"message": "Email atau password salah",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"message": "Email atau password salah",
		})
		return
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":      false,
			"message": "Gagal membuat access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":           true,
		"message":      "Login berhasil",
		"access_token": accessToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}