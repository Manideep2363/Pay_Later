package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"paylater/internal/auth"
	db "paylater/internal/db"
)

type AuthService struct {
	db            *sql.DB
	queries       *db.Queries
	jwtSecret     string
	adminEmail    string
	adminPassword string
}

func NewAuthService(
	dbConn *sql.DB,
	q *db.Queries,
	jwtSecret string,
	adminEmail string,
	adminPassword string,
) *AuthService {

	return &AuthService{
	db:            dbConn,
	queries:       q,
	jwtSecret:     jwtSecret,
	adminEmail:    adminEmail,
	adminPassword: adminPassword,
}}


type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

func (s *AuthService) Register(
	ctx context.Context,
	req RegisterRequest,
) error {

	exists, err := s.queries.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateUser(ctx, db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})

	return err
}

type LoginRequest struct {
	Email    string
	Password string
}


func (s *AuthService) Login(
	ctx context.Context,
	req LoginRequest,
) (string, error) {

	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = auth.CheckPassword(user.Password, req.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(
		user.UserID,
		user.Email,
		"user",
		s.jwtSecret,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

type AdminLoginRequest struct {
	Email    string
	Password string
}


func (s *AuthService) AdminLogin(
	ctx context.Context,
	req AdminLoginRequest,
) (string, error) {

	if req.Email != s.adminEmail || req.Password != s.adminPassword {
		return "", errors.New("invalid admin credentials")
	}

	token, err := auth.GenerateToken(
		0,
		s.adminEmail,
		"admin",
		s.jwtSecret,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

type MerchantRegisterRequest struct {
	Name                   string  `json:"name" binding:"required"`
	Email                  string  `json:"email" binding:"required,email"`
	Phone                  string  `json:"phone" binding:"required"`
	Password               string  `json:"password" binding:"required,min=6"`
	CommissionPercentage   float64 `json:"commission_percentage" binding:"required"`
}

func (s *AuthService) MerchantRegister(
	ctx context.Context,
	req MerchantRegisterRequest,
) error {

	// Hash the merchant password.
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Save the merchant.
	_, err = s.queries.CreateMerchant(ctx, db.CreateMerchantParams{
		Name:                   req.Name,
		Email:                  req.Email,
		Phone:                  req.Phone,
		PasswordHash:           hashedPassword,
		CommissionPercentage:   fmt.Sprintf("%.2f", req.CommissionPercentage),
	})

	return err
}

type MerchantLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) MerchantLogin(
	ctx context.Context,
	req MerchantLoginRequest,
) (string, error) {

	// Get merchant by email
	merchant, err := s.queries.GetMerchantByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare stored hash with the plain password.
	if err := auth.CheckPassword(merchant.PasswordHash, req.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT.
	token, err := auth.GenerateToken(
		merchant.MerchantID,
		merchant.Email,
		"merchant",
		s.jwtSecret,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}