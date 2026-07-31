package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"paylater/internal/config"
	"paylater/internal/db"
	"paylater/internal/handlers"
	"paylater/internal/middleware"
	"paylater/internal/service"
)

func SetupRoutes(
	router *gin.Engine,
	dbConn *sql.DB,
	cfg *config.Config,
) {

	queries := db.New(dbConn)

	// Services
	userService := service.NewUserService(queries)
	merchantService := service.NewMerchantService(queries)
	transactionService := service.NewTransactionService(dbConn, queries)
	reportService := service.NewReportService(queries)
	paymentService := service.NewPaymentService(dbConn, queries)
	authService := service.NewAuthService(
		dbConn,
		queries,
		cfg.JWTSecret,
		cfg.AdminEmail,
		cfg.AdminPassword,
	)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	merchantHandler := handlers.NewMerchantHandler(merchantService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	authHandler := handlers.NewAuthHandler(authService)

	// ===========================
	// Public Routes
	// ===========================
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/merchant/register", authHandler.MerchantRegister)
	router.POST("/merchant/login", authHandler.MerchantLogin)
	router.POST("/admin/login", authHandler.AdminLogin)


	// ===========================
	// User Routes
	// Accessible by:
	// - User
	// - Admin
	// ===========================
	user := router.Group("/")
	user.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
	middleware.RequireRole("user", "admin"),
	)

	merchant := router.Group("/merchant")
	merchant.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
	middleware.RequireRole("merchant", "admin"),
	)

	//merchant.GET("/profile", merchantHandler.GetProfile)

	user.POST("/purchases", transactionHandler.Purchase)
	user.POST("/payments", paymentHandler.Repay)

	user.GET("/users/:id/payments", paymentHandler.ListUserPayments)
	// ===========================
	// Admin Routes
	// Accessible only by Admin
	// ===========================
	admin := router.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RequireRole("admin"),
	)

	// User Management
	admin.POST("/users", userHandler.CreateUser)
	admin.GET("/users", userHandler.ListUsers)

	// Merchant Management
	admin.POST("/merchants", merchantHandler.CreateMerchant)
	admin.GET("/merchants", merchantHandler.ListMerchants)
	admin.GET("/merchants/:id", merchantHandler.GetMerchantByID)
	merchant.GET("/profile", merchantHandler.GetProfile)
	admin.PUT("/merchants/:id/commission", merchantHandler.UpdateMerchantCommission)

	// Transaction Management
	admin.GET("/purchases", transactionHandler.ListTransactions)
	admin.GET("/purchases/:id", transactionHandler.GetTransactionByID)
	admin.GET("/users/:id/purchases", transactionHandler.ListUserTransactions)
	merchant.GET("/merchant/:id/transactions", transactionHandler.ListMerchantTransactions) 

	// Payment Management
	user.GET("/payments", paymentHandler.GetPayments)
	admin.GET("/payments/:id",paymentHandler.GetPaymentByID)

	// Reports
	admin.GET("/reports/outstanding-balance", reportHandler.OutstandingBalance)
	admin.GET("/reports/users-due", reportHandler.UserOutstandingDues)
	admin.GET("/reports/users-at-credit-limit", reportHandler.UsersAtCreditLimit)
	admin.GET("/reports/merchant-commissions", reportHandler.MerchantCommissionSummary)
}