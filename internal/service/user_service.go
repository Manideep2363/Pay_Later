package service

import (
	"context"
	"database/sql"
	"errors"

	"paylater/internal/db"
)

type UserService struct {
	queries *db.Queries
}

// Constructor
func NewUserService(q *db.Queries) *UserService {
	return &UserService{
		queries: q,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	name string,
	email string,
) error {

	_, err := s.queries.GetUserByEmail(ctx, email)

	if err == nil {
		return errors.New("email already exists")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = s.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return err
	}

	return nil
}