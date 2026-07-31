package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"fmt"

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

//Onboarding merchant
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

//Get merchant_By_ID
func (s *MerchantService) GetMerchantByID(
	ctx context.Context,
	id int32,
) (db.GetMerchantByIDRow, error) {

	return s.queries.GetMerchantByID(ctx, id)
}

//ListMerchants
func (s *MerchantService) ListMerchants(
	ctx context.Context) ([]db.ListMerchantsRow,error) {
		return s.queries.ListMerchants(ctx)
}

//update Merchant commission
func (s *MerchantService) UpdateMerchantCommission(
	ctx context.Context,
	id int32,
	commission float64,
) error {

	commissionStr := strconv.FormatFloat(commission, 'f', 2, 64)

	result, err := s.queries.UpdateMerchantCommission(ctx, db.UpdateMerchantCommissionParams{
		MerchantID:             id,
		CommissionPercentage: commissionStr,
	})
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("merchant not found")
	}

	return nil
}