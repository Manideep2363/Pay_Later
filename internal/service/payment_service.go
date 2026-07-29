package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	db "paylater/internal/db"
)

type PaymentService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewPaymentService(dbConn *sql.DB, q *db.Queries) *PaymentService {
	return &PaymentService{
		db:      dbConn,
		queries: q,
	}
}

//PayBack
func (s *PaymentService) Repay(
	ctx context.Context,
	userID int32,
	amount float64,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	user, err := qtx.GetUserByIDForUpdate(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	currentDue, _ := strconv.ParseFloat(user.CurrentDue, 64)

	if amount > currentDue {
		return fmt.Errorf("payment exceeds outstanding due")
	}

	newDue := currentDue - amount
	newDueStr := strconv.FormatFloat(newDue, 'f', 2, 64)
	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)

	err = qtx.CreatePayment(ctx, db.CreatePaymentParams{
		UserID: userID,
		Amount: amountStr,
	})
	if err != nil {
		return err
	}

	_, err = qtx.UpdateUserDue(ctx, db.UpdateUserDueParams{
		UserID:     userID,
		CurrentDue: newDueStr,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

//Payback_History
func (s *PaymentService) GetPaymentByID(
	ctx context.Context,
	id int32,
) (db.Payment, error) {

	return s.queries.GetPaymentByID(ctx, id)
}

//PayBacks made by user
func (s *PaymentService) ListUserPayments(
	ctx context.Context,
	userID int32,
) ([]db.Payment, error) {

	return s.queries.ListUserPayments(ctx, userID)
}