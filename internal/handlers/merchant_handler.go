package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

// MerchantHandler handles all HTTP requests related to merchants.
type MerchantHandler struct {
	// service contains the business logic for merchants.
	service *service.MerchantService
}

// NewMerchantHandler creates and returns a new MerchantHandler.
func NewMerchantHandler(s *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		service: s,
	}
}

// CreateMerchantRequest represents the JSON request body
// required to create a new merchant.
type CreateMerchantRequest struct {
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Commission float64 `json:"commission"`
}

// CreateMerchant handles POST /merchants requests.
func (h *MerchantHandler) CreateMerchant(c *gin.Context) {

	// Variable to store the incoming JSON request.
	var req CreateMerchantRequest

	// Read and bind the JSON request body.
	if err := c.ShouldBindJSON(&req); err != nil {

		// Return 400 Bad Request if the JSON is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service layer to create the merchant.
	err := h.service.CreateMerchant(
		c.Request.Context(),
		req.Name,
		req.Phone,
		req.Commission,
	)

	// Return an error if merchant creation fails.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Merchant created successfully",
	})
}

// GetMerchantByID handles GET /merchants/:id requests.
func (h *MerchantHandler) GetMerchantByID(c *gin.Context) {

	// Read the merchant id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the merchant id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid merchant id",
		})
		return
	}

	// Fetch the merchant from the service layer.
	merchant, err := h.service.GetMerchantByID(
		c.Request.Context(),
		int32(id),
	)

	// Return 404 if the merchant is not found.
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "merchant not found",
		})
		return
	}

	// Return the merchant details.
	c.JSON(http.StatusOK, merchant)
}

func (h *MerchantHandler) GetProfile(c *gin.Context) {

	id := c.MustGet("id").(int32)

	merchant, err := h.service.GetMerchantByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "merchant not found",
		})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// ListMerchants handles GET /merchants requests.
func (h *MerchantHandler) ListMerchants(c *gin.Context) {

	// Fetch all merchants from the service layer.
	merchants, err := h.service.ListMerchants(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the list of merchants.
	c.JSON(http.StatusOK, merchants)
}

// UpdateMerchantCommissionRequest represents the JSON request body
// required to update a merchant's commission.
type UpdateMerchantCommissionRequest struct {
	// Commission must be between 3% and 20%.
	Commission float64 `json:"commission" binding:"required,gte=3,lte=20"`
}

// UpdateMerchantCommission handles
// PUT /merchants/:id/commission requests.
func (h *MerchantHandler) UpdateMerchantCommission(c *gin.Context) {

	// Read the merchant id from the URL.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		// Return 400 if the merchant id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid merchant id",
		})
		return
	}

	// Variable to store the incoming JSON request.
	var req UpdateMerchantCommissionRequest

	// Read and validate the JSON request body.
	if err := c.ShouldBindJSON(&req); err != nil {

		// Return 400 if validation fails.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service layer to update the commission.
	err = h.service.UpdateMerchantCommission(
		c.Request.Context(),
		int32(id),
		req.Commission,
	)

	// Return an error if the merchant is not found
	// or the update fails.
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response.
	c.JSON(http.StatusOK, gin.H{
		"message": "Merchant commission updated successfully",
	})
}