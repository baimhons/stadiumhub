package initial

import (
	middlewareInternal "github.com/baimhons/stadiumhub/internal/middlewares"
	middlewareConfig "github.com/baimhons/stadiumhub/internal/middlewares/configs"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AuthMiddleware middlewareInternal.AuthMiddleware
}

func NewMiddleware(redis utils.RedisClient, jwt utils.JWT, secret string) *Middleware {
	return &Middleware{
		AuthMiddleware: middlewareInternal.NewAuthMiddleware(redis),
	}
}

func SetupMiddleware(app *gin.Engine) {
	app.Use(
		middlewareConfig.CORS(),
	)
}
