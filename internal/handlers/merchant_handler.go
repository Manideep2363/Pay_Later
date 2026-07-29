package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

type MerchantHandler struct {
	service *service.MerchantService
}

func NewMerchantHandler(s *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		service: s,
	}
}

type CreateMerchantRequest struct {
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Commission float64 `json:"commission"`
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {

	var req CreateMerchantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.CreateMerchant(
		c.Request.Context(),
		req.Name,
		req.Phone,
		req.Commission,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Merchant created successfully",
	})
}