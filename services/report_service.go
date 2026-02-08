package services

import (
	"category-api/models"
	"category-api/repositories"
)

type ReportService interface {
	GetDailyReport() (models.DailyReportResponse, error)
}

type reportService struct {
	transactionRepo repositories.TransactionRepository
}

func NewReportService(transactionRepo repositories.TransactionRepository) ReportService {
	return &reportService{transactionRepo}
}

func (s *reportService) GetDailyReport() (models.DailyReportResponse, error) {
	revenue, err := s.transactionRepo.GetTodayRevenue()
	if err != nil {
		return models.DailyReportResponse{}, err
	}

	count, err := s.transactionRepo.GetTodayTransactionCount()
	if err != nil {
		return models.DailyReportResponse{}, err
	}

	bestProduct, err := s.transactionRepo.GetTodayBestSellingProduct()
	if err != nil {
		return models.DailyReportResponse{}, err
	}

	return models.DailyReportResponse{
		TotalRevenue:   revenue,
		TotalTransaksi: count,
		ProdukTerlaris: bestProduct,
	}, nil
}