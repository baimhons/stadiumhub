package user

import (
	"net/http"
	"strings"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
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

	// ตรวจสอบว่าเป็น localhost หรือไม่
	isLocalhost := strings.Contains(c.Request.Host, "localhost") ||
		strings.Contains(c.Request.Host, "127.0.0.1")

	// Set Cookie with dynamic security settings
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   !isLocalhost, // ✅ false สำหรับ localhost, true สำหรับ production
		SameSite: http.SameSiteNoneMode,
	})

	// ส่ง session_id กลับไปใน response body ด้วย (สำหรับ Token-based)
	c.JSON(status, utils.SuccessResponse{
		Message: resp.Message,
		Data: map[string]interface{}{
			"session_id": sessionID, // ✅ เพิ่มบรรทัดนี้!
			"message":    "Cookie has been set (in-memory session)",
			"user":       dataMap["user"], // ส่งข้อมูล user กลับไปด้วย (ถ้ามี)
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

	isLocalhost := strings.Contains(c.Request.Host, "localhost") ||
		strings.Contains(c.Request.Host, "127.0.0.1")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   !isLocalhost,
		SameSite: http.SameSiteNoneMode,
	})

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
