package payment

import (
	"github.com/baimhons/stadiumhub/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type PaymentRoutes struct {
	group           *gin.RouterGroup
	paymentHandler  PaymentHandler
	paymentValidate PaymentValidate
	authMiddleware  middlewares.AuthMiddlewareImpl
}

func NewPaymentRoutes(
	group *gin.RouterGroup,
	paymentHandler PaymentHandler,
	paymentValidate PaymentValidate,
	authMiddleWare middlewares.AuthMiddlewareImpl,
) *PaymentRoutes {

	paymentGroup := group.Group("/payment")
	r := &PaymentRoutes{
		group:           paymentGroup,
		paymentHandler:  paymentHandler,
		paymentValidate: paymentValidate,
		authMiddleware:  authMiddleWare,
	}

	return r
}

func (r *PaymentRoutes) RegisterRoutes() {

	r.group.POST("/create", r.authMiddleware.RequireAuth(), r.paymentValidate.ValidateCreatePayment, r.paymentHandler.CreateCheckoutSession)
}
