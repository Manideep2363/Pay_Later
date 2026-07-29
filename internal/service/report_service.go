package service

import (
	"context"

	"paylater/internal/db"
)

type ReportService struct {
	queries *db.Queries
}

func NewReportService(q *db.Queries) *ReportService {
	return &ReportService{
		queries: q,
	}
}

func (s *ReportService) MerchantCommissionSummary(
	ctx context.Context,
) ([]db.GetMerchantCommissionSummaryRow, error) {

	return s.queries.GetMerchantCommissionSummary(ctx)
}

func (s *ReportService) UserOutstandingDues(
	ctx context.Context,
) ([]db.GetUserOutstandingDuesRow, error) {

	return s.queries.GetUserOutstandingDues(ctx)
}

func (s *ReportService) UsersAtCreditLimit(
	ctx context.Context,
) ([]db.GetUsersAtCreditLimitRow, error) {

	return s.queries.GetUsersAtCreditLimit(ctx)
}

func (s *ReportService) OutstandingBalance(
	ctx context.Context,
) (string, error) {

	return s.queries.GetOutstandingBalance(ctx)
}