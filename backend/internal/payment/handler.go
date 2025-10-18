package payment

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/payment/api/request"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	CreateCheckoutSession(c *gin.Context)
}

type paymentHandlerImpl struct {
	paymentService PaymentService
}

func NewPaymentHandler(paymentService PaymentService) PaymentHandler {
	return &paymentHandlerImpl{
		paymentService: paymentService,
	}
}

func (h *paymentHandlerImpl) CreateCheckoutSession(c *gin.Context) {
	reqRaw, exists := c.Get("req")
	if !exists {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "request not found",
			Error:   nil,
		})
		return
	}

	req, ok := reqRaw.(request.PaymentIntentRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request data",
			Error:   nil,
		})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}
	resp, err := h.paymentService.StripeSession(userCtx.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "make payment success",
		Data:    resp,
	})
}
