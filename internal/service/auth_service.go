package service

import (
	"context"
	"database/sql"
	"errors"

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