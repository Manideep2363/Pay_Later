package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"paylater/internal/db"
	"paylater/internal/handlers"
	"paylater/internal/service"
)

func SetupRoutes(router *gin.Engine, dbConn *sql.DB) {

	queries := db.New(dbConn)

	// Services
	userService := service.NewUserService(queries)
	merchantService := service.NewMerchantService(queries)
	transactionService := service.NewTransactionService(dbConn, queries)
	reportService := service.NewReportService(queries)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	merchantHandler := handlers.NewMerchantHandler(merchantService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	// User Routes
	router.POST("/users", userHandler.CreateUser)

	// Merchant Routes
	router.POST("/merchants", merchantHandler.CreateMerchant)

	// Transaction Routes
	router.POST("/purchase", transactionHandler.Purchase)

	// Report Routes
	router.GET("/reports/outstanding-balance", reportHandler.OutstandingBalance)
	router.GET("/reports/users-due", reportHandler.UserOutstandingDues)
	router.GET("/reports/users-at-credit-limit", reportHandler.UsersAtCreditLimit)
	router.GET("/reports/merchant-commissions", reportHandler.MerchantCommissionSummary)
}