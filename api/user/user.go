package user

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/alan1420/mobile-app-api/repository"

	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	repo *repository.UserRepository
}

func NewUserHandler(
	router *gin.Engine,
	repo *repository.UserRepository,
) *HandlerUser {
	handler := &HandlerUser{repo: repo}
	handler.setupRouter(router)

	return handler
}

func (h *HandlerUser) setupRouter(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.GET("/profile", h.UserProfile)
		user.GET("/tickets", h.UserTickets)
		user.POST("/tickets/create", h.UserCreateTicket)
	}
}

func (h *HandlerUser) UserProfile(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "token is invalid or expired",
		})
		return
	}

	user, statusCode, err := h.repo.GetUser(session.Value)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerUser) UserTickets(c *gin.Context) {
	// verify sessions
	session, err := c.Request.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "token is invalid or expired",
		})
		return
	}

	user, statusCode, err := h.repo.GetUser(session.Value)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	tickets, err := h.repo.GetTicket(user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "ticket not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *HandlerUser) UserCreateTicket(c *gin.Context) {
	// verify sessions
	session, err := c.Request.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "session not found",
		})
		return
	}

	user, statusCode, err := h.repo.GetUser(session.Value)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	statusCode, err = h.repo.CreateTicket(user.ID, data)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ticket created",
	})
}
