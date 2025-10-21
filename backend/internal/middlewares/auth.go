package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/utils"

	"github.com/gin-gonic/gin"
)

// ---- In-memory session store ----
var (
	sessionStore = make(map[string]string)
	sessionMutex sync.RWMutex
)

// ฟังก์ชันช่วย set/get session
func SetSession(sessionID string, userCtx models.UserContext, exp time.Duration) {
	userCtxJSON, _ := json.Marshal(userCtx)
	sessionMutex.Lock()
	sessionStore[sessionID] = string(userCtxJSON)
	sessionMutex.Unlock()

	// ตั้งเวลาหมดอายุ
	go func() {
		time.Sleep(exp)
		sessionMutex.Lock()
		delete(sessionStore, sessionID)
		sessionMutex.Unlock()
	}()
}

func getSession(sessionID string) (string, bool) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	data, ok := sessionStore[sessionID]
	return data, ok
}

func DeleteSession(sessionID string) {
	sessionMutex.Lock()
	delete(sessionStore, sessionID)
	sessionMutex.Unlock()
}

// ---- Middleware RequireAuth ----
type AuthMiddlewareImpl struct{}

func (a *AuthMiddlewareImpl) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "session not found",
				Error:   err,
			})
			return
		}

		// ดึงข้อมูลจาก in-memory cache
		userContextJSON, ok := getSession(sessionID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "session expired or invalid",
				Error:   fmt.Errorf("session not found in memory"),
			})
			return
		}

		var userContext models.UserContext
		if err := json.Unmarshal([]byte(userContextJSON), &userContext); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unmarshal error: %v", err)})
			return
		}

		c.Set("userContext", userContext)
		c.Next()
	}
}
