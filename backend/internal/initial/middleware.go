package initial

import (
	middlewareInternal "github.com/baimhons/stadiumhub/internal/middlewares"
	middlewareConfig "github.com/baimhons/stadiumhub/internal/middlewares/configs"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AuthMiddleware middlewareInternal.AuthMiddlewareImpl
}

func NewMiddleware() *Middleware {
	return &Middleware{
		AuthMiddleware: middlewareInternal.AuthMiddlewareImpl{},
	}
}

func SetupMiddleware(app *gin.Engine) {
	app.Use(
		middlewareConfig.CORS(),
	)
}
