package user

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type userHandlerImpl struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &userHandlerImpl{userService: userService}
}

func (h *userHandlerImpl) RegisterUser(c *gin.Context) {
	req, ok := c.Get("req")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	registerReq, ok := req.(request.RegisterUser)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request type"})
		return
	}

	resp, status, err := h.userService.RegisterUser(registerReq)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, resp)
}

func (h *userHandlerImpl) LoginUser(c *gin.Context) {
	req, ok := c.Get("req")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	loginReq, ok := req.(request.LoginUser)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request type"})
		return
	}

	resp, status, err := h.userService.LoginUser(loginReq)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, resp)
}
