package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

type PurchaseRequest struct {
	UserID     int32   `json:"user_id"`
	MerchantID int32   `json:"merchant_id"`
	Amount     float64 `json:"amount"`
}

func (h *TransactionHandler) Purchase(c *gin.Context) {

	var req PurchaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Purchase(
		c.Request.Context(),
		req.UserID,
		req.MerchantID,
		req.Amount,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Purchase successful",
	})
}