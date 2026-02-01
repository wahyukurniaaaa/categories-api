package services

import (
	"category-api/models"
	"category-api/repositories"

	"github.com/google/uuid"
)

type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id string) (models.Category, error)
	Create(c models.Category) (models.Category, error)
	Update(id string, c models.Category) (models.Category, error)
	Delete(id string) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id string) (models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Create(c models.Category) (models.Category, error) {
	// Generate UUID here if not present
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return s.repo.Create(c)
}

func (s *categoryService) Update(id string, c models.Category) (models.Category, error) {
	return s.repo.Update(id, c)
}

func (s *categoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
