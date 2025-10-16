package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
	redisLib "github.com/redis/go-redis/v9"
)

type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	redis utils.RedisClient
}

func NewAuthMiddleware(redis utils.RedisClient) *AuthMiddlewareImpl {
	return &AuthMiddlewareImpl{redis: redis}
}

func (a *AuthMiddlewareImpl) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// อ่าน session ID จาก cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "session not found",
				Error:   err,
			})
			return
		}

		// ดึง user context จาก Redis
		userContextJSON, err := a.redis.Get(
			context.Background(),
			fmt.Sprintf("session:%s", sessionID),
		)
		if err != nil {
			if err == redisLib.Nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
					Message: "session expired or invalid",
					Error:   err,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("redis error: %v", err)})
			return
		}

		var userContext models.UserContext
		if err := json.Unmarshal([]byte(userContextJSON), &userContext); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unmarshal error: %v", err)})
			return
		}

		// เก็บ user context ไว้ใน context
		c.Set("userContext", userContext)
		c.Next()
	}
}
