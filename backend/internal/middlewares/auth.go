package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
	redisLib "github.com/redis/go-redis/v9"
)

type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	redis  utils.RedisClient
	jwt    utils.JWT
	secret string
}

func NewAuthMiddleware(redis utils.RedisClient, jwt utils.JWT, secret string) *AuthMiddlewareImpl {
	return &AuthMiddlewareImpl{redis: redis, jwt: jwt, secret: secret}
}

func (a *AuthMiddlewareImpl) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "[jwt] invalid token type"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		var tokenContext user.TokenContext
		_, err := a.jwt.Parse(token, &tokenContext, a.secret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("[jwt] %v", err)})
			return
		}

		userContextJSON, err := a.redis.Get(context.Background(), fmt.Sprintf("access_token:%s", tokenContext.UserID))
		if err != nil {
			if err == redisLib.Nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "[jwt] session not found"})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("[jwt] %v", err)})
			return
		}

		var userContext user.UserContext
		if err := json.Unmarshal([]byte(userContextJSON), &userContext); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("[jwt] %v", err)})
			return
		}

		c.Set("userContext", userContext)
		c.Next()
	}
}
