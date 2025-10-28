package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func SetSession(sessionID string, userCtx models.UserContext, exp time.Duration) {
	userCtxJSON, _ := json.Marshal(userCtx)
	sessionMutex.Lock()
	sessionStore[sessionID] = string(userCtxJSON)
	sessionMutex.Unlock()

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
		var sessionID string
		var signature string

		// 1. ลองอ่านจาก Cookie ก่อน (สำหรับ localhost)
		cookieValue, err := c.Cookie("session_id")
		if err == nil && cookieValue != "" {
			// แยก sessionID|signature
			parts := strings.Split(cookieValue, "|")
			if len(parts) == 2 {
				sessionID = parts[0]
				signature = parts[1]
			} else {
				sessionID = cookieValue // ถ้าไม่มี signature
			}
		}

		// 2. ถ้าไม่มี Cookie ลองอ่านจาก Authorization header (สำหรับ production)
		if sessionID == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// ลบ "Bearer " prefix
				token := strings.TrimPrefix(authHeader, "Bearer ")
				token = strings.TrimSpace(token)

				// แยก sessionID|signature (ถ้ามี)
				parts := strings.Split(token, "|")
				if len(parts) == 2 {
					sessionID = parts[0]
					signature = parts[1]
				} else {
					sessionID = token
				}
			}
		}

		// 3. ถ้ายังไม่มี session_id = Unauthorized
		if sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "session not found",
				Error:   fmt.Errorf("no session_id in cookie or authorization header"),
			})
			return
		}

		// 4. ตรวจสอบลายเซ็น (ถ้ามี)
		if signature != "" && !utils.VerifySession(sessionID, signature) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "cookie signature invalid",
			})
			return
		}

		// 5. ดึงข้อมูลจาก in-memory cache
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("unmarshal error: %v", err),
			})
			return
		}

		c.Set("userContext", userContext)
		c.Next()
	}
}
