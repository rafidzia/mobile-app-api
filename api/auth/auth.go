package auth

import (
	"io"
	"net/http"

	"github.com/alan1420/mobile-app-api/repository"

	"github.com/gin-gonic/gin"
)

type HandlerAuth struct {
	repo *repository.AuthRepository
}

func NewAuthHandler(
	router *gin.Engine,
	repo *repository.AuthRepository,
) *HandlerAuth {
	handler := &HandlerAuth{
		repo: repo,
	}
	handler.setupRouter(router)

	return handler
}

func (h *HandlerAuth) setupRouter(router *gin.Engine) {
	router.POST("/login", h.UserLogin)
	router.POST("/register", h.UserRegister)
	router.POST("/forgot-password", h.UserForgotPassword)
}

func (h *HandlerAuth) UserLogin(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	sessionId, statusCode, err := h.repo.AuthLoginUser(data)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.SetCookie(
		"session_id",
		sessionId,
		3600,
		"/",
		"localhost:8080",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *HandlerAuth) UserRegister(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	statusCode, err := h.repo.AuthRegisterUser(data)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *HandlerAuth) UserForgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "no implementation yet",
	})
}
