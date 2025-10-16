package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/baimhons/stadiumhub/internal/user/api/response"
	"github.com/baimhons/stadiumhub/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(req request.RegisterUser) (resp utils.SuccessResponse, statusCode int, err error)
	LoginUser(req request.LoginUser) (resp utils.SuccessResponse, statusCode int, err error)
	LogoutUser(userCtx models.UserContext, sessionID string) (int, error)
	GetUserProfile(userCtx models.UserContext) (resp utils.SuccessResponse, statusCode int, err error)
	UpdateUser(userCtx models.UserContext, req request.UpdateUser) (resp utils.SuccessResponse, statusCode int, err error)
}

type userServiceImpl struct {
	userRepository UserRepository
	redis          utils.RedisClient
}

func NewUserService(userRepository UserRepository, redis utils.RedisClient) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		redis:          redis,
	}
}

func (us *userServiceImpl) RegisterUser(req request.RegisterUser) (resp utils.SuccessResponse, statusCode int, err error) {
	user := User{}

	if err := us.userRepository.GetBy("email", req.Email, &user); err == nil {
		return resp, http.StatusConflict, errors.New("email already exists")
	}
	if err := us.userRepository.GetBy("username", req.Username, &user); err == nil {
		return resp, http.StatusConflict, errors.New("username already exists")
	}

	newUser := User{
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        "user",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp, http.StatusInternalServerError, err
	}
	newUser.Password = string(hashPassword)

	tx := us.userRepository.Begin()
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return resp, http.StatusInternalServerError, err
	}
	tx.Commit()

	return utils.SuccessResponse{
		Message: "User registered successfully!",
		Data:    nil,
	}, http.StatusCreated, nil
}

func (us *userServiceImpl) LoginUser(req request.LoginUser) (resp utils.SuccessResponse, statusCode int, err error) {
	user := User{}

	if err := us.userRepository.GetByUsernameOrEmail(req.UsernameOrEmail, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, http.StatusNotFound, errors.New("user not found")
		}
		return resp, http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return resp, http.StatusUnauthorized, errors.New("invalid credentials")
	}

	// สร้าง session ID แบบ random
	sessionID, err := utils.GenerateSecureToken(32)
	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	sessionExp := time.Hour * 24 // session หมดอายุใน 24 ชั่วโมง

	// เก็บ user context ใน Redis โดยใช้ session ID เป็น key
	userContext := models.UserContext{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
	}

	userContextJSON, _ := json.Marshal(userContext)

	// เก็บ session ใน Redis
	if err := us.redis.Set(
		context.Background(),
		fmt.Sprintf("session:%s", sessionID),
		userContextJSON,
		sessionExp,
	); err != nil {
		return resp, http.StatusInternalServerError, err
	}

	// return sessionID เพื่อให้ handler ไปตั้งค่า cookie
	// แต่ไม่ส่งไปที่ client ใน response
	return utils.SuccessResponse{
		Message: "User logged in successfully!",
		Data: map[string]interface{}{
			"sessionID": sessionID, // ใช้ภายในเท่านั้น
		},
	}, http.StatusOK, nil
}

// LogoutUser - ลบ session (แก้ไข signature)
func (us *userServiceImpl) LogoutUser(userCtx models.UserContext, sessionID string) (int, error) {
	err := us.redis.Del(
		context.Background(),
		fmt.Sprintf("session:%s", sessionID),
	)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (us *userServiceImpl) GetUserProfile(userCtx models.UserContext) (resp utils.SuccessResponse, statusCode int, err error) {
	user := User{}
	if err := us.userRepository.GetBy("id", userCtx.ID.String(), &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, http.StatusNotFound, errors.New("user not found")
		}
		return resp, http.StatusInternalServerError, err
	}

	return utils.SuccessResponse{
		Message: "User profile fetched successfully!",
		Data: response.UserProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
		},
	}, http.StatusOK, nil
}

func (us *userServiceImpl) UpdateUser(userCtx models.UserContext, req request.UpdateUser) (resp utils.SuccessResponse, statusCode int, err error) {
	user := User{}

	if err := us.userRepository.GetBy("id", userCtx.ID.String(), &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, http.StatusNotFound, errors.New("user not found")
		}
		return resp, http.StatusInternalServerError, err
	}

	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.PhoneNumber = req.PhoneNumber

	if err := us.userRepository.Update(&user); err != nil {
		return resp, http.StatusInternalServerError, err
	}

	return utils.SuccessResponse{
		Message: "User updated successfully!",
		Data: response.UserProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
		},
	}, http.StatusOK, nil
}
