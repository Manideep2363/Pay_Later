package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{
		service: s,
	}
}

func (h *ReportHandler) MerchantCommissionSummary(c *gin.Context) {

	result, err := h.service.MerchantCommissionSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReportHandler) UserOutstandingDues(c *gin.Context) {

	result, err := h.service.UserOutstandingDues(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReportHandler) UsersAtCreditLimit(c *gin.Context) {

	result, err := h.service.UsersAtCreditLimit(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReportHandler) OutstandingBalance(c *gin.Context) {

	result, err := h.service.OutstandingBalance(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_outstanding_balance": result,
	})
}