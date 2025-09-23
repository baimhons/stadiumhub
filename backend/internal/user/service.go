package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/user/api/request"
	"github.com/baimhons/stadiumhub/internal/user/api/response"
	"github.com/baimhons/stadiumhub/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	LoginUser(req request.LoginUser) (resp utils.SuccessResponse, statusCode int, err error)
	RegisterUser(req request.RegisterUser) (resp utils.SuccessResponse, statusCode int, err error)
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
		FullName:    req.FullName,
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
	// หาด้วย email หรือ username
	if err := us.userRepository.GetByUsernameOrEmail(req.UsernameOrEmail, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, http.StatusNotFound, errors.New("user not found")
		}
		return resp, http.StatusInternalServerError, err
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return resp, http.StatusUnauthorized, errors.New("invalid credentials")
	}

	// ===== Generate Token =====
	timeNow := time.Now()
	accessTokenExp := timeNow.Add(time.Hour * 1)
	refreshTokenExp := timeNow.Add(time.Hour * 24)
	secret := internal.ENV.JWTSecret.Secret

	accessToken, err := utils.NewJWT().Generate(map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	}, accessTokenExp.Unix(), secret)

	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	refreshToken, err := utils.NewJWT().Generate(map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	}, refreshTokenExp.Unix(), secret)
	if err != nil {
		return resp, http.StatusInternalServerError, err
	}

	// เก็บ token ไว้ใน Redis
	if err := us.redis.Set(context.Background(), fmt.Sprintf("access_token:%s", user.ID), accessToken, accessTokenExp.Sub(timeNow)); err != nil {
		return resp, http.StatusInternalServerError, err
	}
	if err := us.redis.Set(context.Background(), fmt.Sprintf("refresh_token:%s", user.ID), refreshToken, refreshTokenExp.Sub(timeNow)); err != nil {
		return resp, http.StatusInternalServerError, err
	}

	return utils.SuccessResponse{
		Message: "User logged in successfully!",
		Data: response.LoginUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, http.StatusOK, nil
}
