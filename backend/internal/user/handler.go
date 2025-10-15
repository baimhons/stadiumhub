package user

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	LogoutUser(c *gin.Context)
	GetUserProfile(c *gin.Context)
	UpdateUser(c *gin.Context)
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

func (h *userHandlerImpl) LogoutUser(c *gin.Context) {
	userCtx, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	ctx, ok := userCtx.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user context"})
		return
	}

	status, err := h.userService.LogoutUser(ctx)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, gin.H{"message": "User logged out successfully"})
}

func (h *userHandlerImpl) GetUserProfile(c *gin.Context) {
	userCtx, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	ctx, ok := userCtx.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user context"})
		return
	}

	resp, status, err := h.userService.GetUserProfile(ctx)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, resp)
}

func (h *userHandlerImpl) UpdateUser(c *gin.Context) {
	req, exists := c.Get("req")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	updateReq, ok := req.(request.UpdateUser)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request type"})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user context"})
		return
	}

	resp, status, err := h.userService.UpdateUser(userCtx, updateReq)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, resp)
}
