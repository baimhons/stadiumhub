package user

import (
	"fmt"
	"net/http"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type userValidate struct {
}

type UserValidate interface {
	ValidateRegisterUser(c *gin.Context)
	ValidateLoginUser(c *gin.Context)
	ValidateUpdateUser(c *gin.Context)
	ValidateRoleAdmin(c *gin.Context)
}

func NewUserValidate() *userValidate {
	return &userValidate{}
}

func (u *userValidate) ValidateRegisterUser(c *gin.Context) {
	var req request.RegisterUser

	if err := utils.ValidateCommonRequestBody(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	if req.Username == "admin" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "username cannot be admin",
		})
		c.Abort()
		return
	}

	c.Set("req", req)
	c.Next()
}

func (u *userValidate) ValidateLoginUser(c *gin.Context) {
	var req request.LoginUser

	if err := utils.ValidateCommonRequestBody(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	c.Set("req", req)
	c.Next()
}

func (u *userValidate) ValidateUpdateUser(c *gin.Context) {
	var req request.UpdateUser

	if err := utils.ValidateCommonRequestBody(c, &req); err != nil {
		fmt.Println("err")
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	c.Set("req", req)
	c.Next()
}

func (u *userValidate) ValidateRoleAdmin(c *gin.Context) {
	userContext, ok := c.Get("userContext")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "Unauthorized",
		})
		c.Abort()
		return
	}

	if uc, ok := userContext.(models.UserContext); ok {
		if uc.Role != "admin" {
			c.JSON(http.StatusForbidden, utils.ErrorResponse{
				Message: "Forbidden",
			})
			c.Abort()
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "Invalid user context",
		})
		c.Abort()
		return
	}

	c.Next()
}
