package payment

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/payment/api/request"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

// validate payment request
type paymentValidate struct {
}

type PaymentValidate interface {
	ValidateCreatePayment(c *gin.Context)
}

func NewPaymentValidate() *paymentValidate {
	return &paymentValidate{}
}

func (p *paymentValidate) ValidateCreatePayment(c *gin.Context) {
	var req request.PaymentIntentRequest

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
