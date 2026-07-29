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
	paymentService := service.NewPaymentService(dbConn, queries)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	merchantHandler := handlers.NewMerchantHandler(merchantService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// User Routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.ListUsers)

	// Merchant Routes
	router.POST("/merchants", merchantHandler.CreateMerchant)
	router.GET("/merchants/:id", merchantHandler.GetMerchantByID)
	router.GET("/merchants",merchantHandler.ListMerchants)
	router.PUT("/merchants/:id/commission", merchantHandler.UpdateMerchantCommission)

	// Transaction Routes
	router.POST("/purchase", transactionHandler.Purchase)
	router.GET("/transactions", transactionHandler.ListTransactions)
	router.GET("/transactions/:id", transactionHandler.GetTransactionByID)
	router.GET("/users/:id/transactions", transactionHandler.ListUserTransactions)

	// Report Routes
	router.GET("/reports/outstanding-balance", reportHandler.OutstandingBalance)
	router.GET("/reports/users-due", reportHandler.UserOutstandingDues)
	router.GET("/reports/users-at-credit-limit", reportHandler.UsersAtCreditLimit)
	router.GET("/reports/merchant-commissions", reportHandler.MerchantCommissionSummary)

	//paybacks Routes
	router.POST("/repay", paymentHandler.Repay)
	router.GET("/payments/:id", paymentHandler.GetPaymentByID)
	router.GET("/users/:id/payments", paymentHandler.ListUserPayments)
}