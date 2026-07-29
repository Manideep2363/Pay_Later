package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

// ReportHandler handles all HTTP requests related to reports.
type ReportHandler struct {
	// service contains the business logic for generating reports.
	service *service.ReportService
}

// NewReportHandler creates and returns a new ReportHandler.
func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{
		service: s,
	}
}

// MerchantCommissionSummary handles
// GET /reports/merchant-commissions requests.
func (h *ReportHandler) MerchantCommissionSummary(c *gin.Context) {

	// Fetch the merchant commission summary from the service layer.
	result, err := h.service.MerchantCommissionSummary(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the commission summary.
	c.JSON(http.StatusOK, result)
}

// UserOutstandingDues handles
// GET /reports/user-outstanding-dues requests.
func (h *ReportHandler) UserOutstandingDues(c *gin.Context) {

	// Fetch the outstanding dues for all users.
	result, err := h.service.UserOutstandingDues(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the outstanding dues report.
	c.JSON(http.StatusOK, result)
}

// UsersAtCreditLimit handles
// GET /reports/users-at-credit-limit requests.
func (h *ReportHandler) UsersAtCreditLimit(c *gin.Context) {

	// Fetch users who have reached their credit limit.
	result, err := h.service.UsersAtCreditLimit(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the list of users at their credit limit.
	c.JSON(http.StatusOK, result)
}

// OutstandingBalance handles
// GET /reports/outstanding-balance requests.
func (h *ReportHandler) OutstandingBalance(c *gin.Context) {

	// Calculate the total outstanding balance.
	result, err := h.service.OutstandingBalance(c.Request.Context())

	// Return 500 if something goes wrong.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the total outstanding balance.
	c.JSON(http.StatusOK, gin.H{
		"total_outstanding_balance": result,
	})
}