package repositories

import (
	"category-api/models"
	"database/sql"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id string) (models.Category, error)
	Create(c models.Category) (models.Category, error)
	Update(id string, c models.Category) (models.Category, error)
	Delete(id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(id string) (models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.Description)
	return c, err
}

func (r *categoryRepository) Create(c models.Category) (models.Category, error) {
	// Assuming ID is passed from service (UUID) or generated here.
	// The implementation plan says "Create Category" gets JSON.
	// Since usage is uuid, we will insert the ID provided by struct.
	_, err := r.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
		c.ID, c.Name, c.Description)
	return c, err
}

func (r *categoryRepository) Update(id string, c models.Category) (models.Category, error) {
	_, err := r.db.Exec("UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		c.Name, c.Description, id)
	if err != nil {
		return c, err
	}
	c.ID = id
	return c, nil
}

func (r *categoryRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
