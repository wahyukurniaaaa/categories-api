package services

import (
	"category-api/models"
	"category-api/repositories"
)

type ProductService interface {
	GetAll() ([]models.Produk, error)
	GetByID(id int) (models.Produk, error)
	Create(p models.Produk) (models.Produk, error)
	Update(id int, p models.Produk) (models.Produk, error)
	Delete(id int) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAll() ([]models.Produk, error) {
	return s.repo.GetAll()
}

func (s *productService) GetByID(id int) (models.Produk, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Create(p models.Produk) (models.Produk, error) {
	return s.repo.Create(p)
}

func (s *productService) Update(id int, p models.Produk) (models.Produk, error) {
	return s.repo.Update(id, p)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
