package services

import (
	"category-api/models"
	"category-api/repositories"
	"errors"
)

type TransactionService interface {
	Checkout(req models.CheckoutRequest) (models.CheckoutResponse, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, productRepo repositories.ProductRepository) TransactionService {
	return &transactionService{transactionRepo, productRepo}
}

func (s *transactionService) Checkout(req models.CheckoutRequest) (models.CheckoutResponse, error) {
	if len(req.Items) == 0 {
		return models.CheckoutResponse{}, errors.New("items cannot be empty")
	}

	// Validate products and stock
	products := make(map[int]models.Produk)
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return models.CheckoutResponse{}, errors.New("quantity must be greater than 0")
		}

		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return models.CheckoutResponse{}, errors.New("product not found")
		}

		if product.Stok < item.Quantity {
			return models.CheckoutResponse{}, errors.New("insufficient stock for product: " + product.Nama)
		}

		products[item.ProductID] = product
	}

	// Create transaction
	transaction, details, err := s.transactionRepo.CreateTransaction(req.Items, products)
	if err != nil {
		return models.CheckoutResponse{}, err
	}

	return models.CheckoutResponse{
		Transaction: transaction,
		Details:     details,
	}, nil
}