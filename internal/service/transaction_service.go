package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"paylater/internal/db"
)

type TransactionService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewTransactionService(dbConn *sql.DB, q *db.Queries) *TransactionService {
	return &TransactionService{
		db:      dbConn,
		queries: q,
	}
}

func (s *TransactionService) Purchase(
	ctx context.Context,
	userID int32,
	merchantID int32,
	amount float64,
) error {

	// Business Validation
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Begin Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	// Lock user row
	user, err := qtx.GetUserByIDForUpdate(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	// Get merchant
	merchant, err := qtx.GetMerchantByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("merchant not found")
		}
		return err
	}

	// Convert DECIMAL strings to float64
	creditLimit, err := parseDecimal(user.CreditLimit)
	if err != nil {
		return err
	}

	currentDue, err := parseDecimal(user.CurrentDue)
	if err != nil {
		return err
	}

	commissionPercentage, err := parseDecimal(merchant.CommissionPercentage)
	if err != nil {
		return err
	}

	// Check available credit
	availableCredit := creditLimit - currentDue

	if amount > availableCredit {
		return errors.New("insufficient credit limit")
	}

	// Calculate commission
	commissionAmount := amount * commissionPercentage / 100

	// Update user's outstanding due
	newDue := currentDue + amount

	_, err = qtx.UpdateUserDue(ctx, db.UpdateUserDueParams{
		UserID:     user.UserID,
		CurrentDue: formatDecimal(newDue),
	})
	if err != nil {
		return err
	}

	// Create transaction
	_, err = qtx.CreateTransaction(ctx, db.CreateTransactionParams{
		UserID:               user.UserID,
		MerchantID:           merchant.MerchantID,
		Amount:               formatDecimal(amount),
		CommissionPercentage: formatDecimal(commissionPercentage),
		CommissionAmount:     formatDecimal(commissionAmount),
	})
	if err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func parseDecimal(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func formatDecimal(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

//List_out_all_Transactions
func (s *TransactionService) ListTransactions(
	ctx context.Context,
) ([]db.Transaction, error) {

	return s.queries.ListTransactions(ctx)
}

//Access Tranasaction By ID
func (s *TransactionService) GetTransactionByID(
	ctx context.Context,
	id int32,
) (db.Transaction, error) {

	return s.queries.GetTransactionByID(ctx, id)
}

//Transactions made by single user
func (s *TransactionService) ListUserTransactions(
	ctx context.Context,
	userID int32,
) ([]db.Transaction, error) {

	return s.queries.ListUserTransactions(ctx, userID)
}

//Transactions of single merchant
func (s *TransactionService) ListMerchantTransactions(
	ctx context.Context,
	merchantID int32,
) ([]db.Transaction, error) {

	return s.queries.ListMerchantTransactions(ctx, merchantID)
}