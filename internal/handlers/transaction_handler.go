package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

// TransactionHandler handles all HTTP requests related to transactions.
type TransactionHandler struct {
	// service contains the business logic for transactions.
	service *service.TransactionService
}

// NewTransactionHandler creates and returns a new TransactionHandler.
func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

// PurchaseRequest represents the JSON request body
// required to make a purchase.
type PurchaseRequest struct {
	MerchantID int32   `json:"merchant_id"`
	Amount     float64 `json:"amount"`
}

// Purchase handles POST /transactions/purchase requests.
func (h *TransactionHandler) Purchase(c *gin.Context) {

	// Variable to store the incoming JSON request.
	var req PurchaseRequest

	// Read and bind the JSON request body.
	if err := c.ShouldBindJSON(&req); err != nil {

		// Return 400 Bad Request if the JSON is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service layer to process the purchase.

	userID := c.MustGet("id").(int32)
	err := h.service.Purchase(
		c.Request.Context(),
		userID,
		req.MerchantID,
		req.Amount,
	)

	// Return an error if the purchase fails.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Purchase successful",
	})
}

// ListTransactions handles GET /transactions requests.
func (h *TransactionHandler) ListTransactions(c *gin.Context) {

	// Fetch all transactions from the service layer.
	transactions, err := h.service.ListTransactions(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the list of transactions.
	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID handles GET /transactions/:id requests.
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {

	// Read the transaction id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction id",
		})
		return
	}

	// Fetch the transaction from the service layer.
	transaction, err := h.service.GetTransactionByID(
		c.Request.Context(),
		int32(id),
	)

	// Return 404 if the transaction is not found.
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	// Return the transaction.
	c.JSON(http.StatusOK, transaction)
}

// ListUserTransactions handles
// GET /users/:id/transactions requests.
func (h *TransactionHandler) ListUserTransactions(c *gin.Context) {

	// Read the user id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the user id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// Fetch all transactions for the given user.
	transactions, err := h.service.ListUserTransactions(
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

	// Return the user's transactions.
	c.JSON(http.StatusOK, transactions)
}

// ListMerchantTransactions handles
// GET /merchants/:id/transactions requests.
func (h *TransactionHandler) ListMerchantTransactions(c *gin.Context) {

	// Read the merchant id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the merchant id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid merchant id",
		})
		return
	}

	// Fetch all transactions for the given merchant.
	transactions, err := h.service.ListMerchantTransactions(
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

	// Return the merchant's transactions.
	c.JSON(http.StatusOK, transactions)
}