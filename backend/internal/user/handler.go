package user

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/baimhons/stadiumhub/internal/user/api/response"
	"github.com/baimhons/stadiumhub/internal/utils"
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
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request",
			Error:   nil,
		})
		return
	}

	registerReq, ok := req.(request.RegisterUser)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request type",
			Error:   nil,
		})
		return
	}

	resp, status, err := h.userService.RegisterUser(registerReq)
	if err != nil {
		c.JSON(status, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(status, resp)
}

func (h *userHandlerImpl) LoginUser(c *gin.Context) {
	req, ok := c.Get("req")
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Message: "invalid request"})
		return
	}

	loginReq, ok := req.(request.LoginUser)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Message: "invalid request type"})
		return
	}

	resp, status, err := h.userService.LoginUser(loginReq)
	if err != nil {
		c.JSON(status, utils.ErrorResponse{Message: err.Error(), Error: err})
		return
	}

	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Message: "invalid response data"})
		return
	}

	sessionID, ok := dataMap["session_id"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Message: "invalid session ID"})
		return
	}

	c.SetCookie("session_id", sessionID, 86400, "/", "", false, true)

	c.JSON(status, utils.SuccessResponse{
		Message: resp.Message,
		Data: response.LoginUserResponse{
			Message: "Cookie has been set (in-memory session)",
		},
	})
}

func (h *userHandlerImpl) LogoutUser(c *gin.Context) {
	userCtx, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}

	ctx, ok := userCtx.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "session not found",
			Error:   err,
		})
		return
	}

	status, err := h.userService.LogoutUser(ctx, sessionID)
	if err != nil {
		c.JSON(status, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	// ลบ cookie
	c.SetCookie(
		"session_id",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(status, utils.SuccessResponse{
		Message: "User logged out successfully",
		Data:    nil,
	})
}

func (h *userHandlerImpl) GetUserProfile(c *gin.Context) {
	userCtx, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}

	ctx, ok := userCtx.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	resp, status, err := h.userService.GetUserProfile(ctx)
	if err != nil {
		c.JSON(status, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(status, resp)
}

func (h *userHandlerImpl) UpdateUser(c *gin.Context) {
	req, exists := c.Get("req")
	if !exists {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request",
			Error:   nil,
		})
		return
	}

	updateReq, ok := req.(request.UpdateUser)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request type",
			Error:   nil,
		})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}

	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	resp, status, err := h.userService.UpdateUser(userCtx, updateReq)
	if err != nil {
		c.JSON(status, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(status, resp)
}
