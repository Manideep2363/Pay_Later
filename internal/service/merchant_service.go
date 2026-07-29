package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"paylater/internal/db"
)

type MerchantService struct {
	queries *db.Queries
}

func NewMerchantService(q *db.Queries) *MerchantService {
	return &MerchantService{
		queries: q,
	}
}

func (s *MerchantService) CreateMerchant(
	ctx context.Context,
	name string,
	phone string,
	commission float64,
) error {

	_, err := s.queries.GetMerchantByPhone(ctx, phone)

	if err == nil {
		return errors.New("merchant phone already exists")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	commissionStr := strconv.FormatFloat(commission, 'f', 2, 64)

	_, err = s.queries.CreateMerchant(ctx, db.CreateMerchantParams{
		Name:                 name,
		Phone:                phone,
		CommissionPercentage: commissionStr,
	})

	if err != nil {
		return err
	}

	return nil
}