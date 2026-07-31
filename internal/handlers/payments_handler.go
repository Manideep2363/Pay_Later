package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

// PaymentHandler handles all HTTP requests related to payments.
type PaymentHandler struct {
	// service contains the business logic for payments.
	service *service.PaymentService
}

// NewPaymentHandler creates and returns a new PaymentHandler.
func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: s,
	}
}

// RepayRequest represents the JSON request body
// required to make a payment.
type RepayRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// Repay handles POST /payments/repay requests.
func (h *PaymentHandler) Repay(c *gin.Context) {

	// Variable to store the incoming JSON request.
	var req RepayRequest

	// Read and bind the JSON request body.
	if err := c.ShouldBindJSON(&req); err != nil {

		// Return 400 Bad Request if the JSON is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service layer to process the payment.
	userID := c.MustGet("id").(int32)
	err := h.service.Repay(
		c.Request.Context(),
		userID,
		req.Amount,
	)

	// Return an error if the payment fails.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response.
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful",
	})
}

// GetPaymentByID handles GET /payments/:id requests.
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {

	// Read the payment id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the payment id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payment id",
		})
		return
	}

	// Fetch the payment from the service layer.
	payment, err := h.service.GetPaymentByID(
		c.Request.Context(),
		int32(id),
	)

	// Return 404 if the payment is not found.
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "payment not found",
		})
		return
	}

	// Return the payment details.
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPayments(c *gin.Context) {

	userID := c.MustGet("id").(int32)

	payments, err := h.service.ListUserPayments(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ListUserPayments handles
// GET /users/:id/payments requests.
func (h *PaymentHandler) ListUserPayments(c *gin.Context) {

	// Read the user id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the user id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// Fetch all payments made by the given user.
	payments, err := h.service.ListUserPayments(
		c.Request.Context(),
		int32(id),
	)

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the user's payment history.
	c.JSON(http.StatusOK, payments)
}