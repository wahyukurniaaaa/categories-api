package repositories

import (
	"category-api/models"
	"database/sql"
)

type ProductRepository interface {
	GetAll(name string) ([]models.Produk, error)
	GetByID(id int) (models.Produk, error)
	Create(p models.Produk) (models.Produk, error)
	Update(id int, p models.Produk) (models.Produk, error)
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAll(name string) ([]models.Produk, error) {
	var rows *sql.Rows
	var err error

	if name != "" {
		// Search by name using ILIKE for case-insensitive matching
		rows, err = r.db.Query("SELECT id, nama, harga, stok FROM products WHERE nama ILIKE '%' || $1 || '%'", name)
	} else {
		rows, err = r.db.Query("SELECT id, nama, harga, stok FROM products")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Produk
	for rows.Next() {
		var p models.Produk
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) GetByID(id int) (models.Produk, error) {
	var p models.Produk
	err := r.db.QueryRow("SELECT id, nama, harga, stok FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)
	return p, err
}

func (r *productRepository) Create(p models.Produk) (models.Produk, error) {
	err := r.db.QueryRow("INSERT INTO products (nama, harga, stok) VALUES ($1, $2, $3) RETURNING id",
		p.Nama, p.Harga, p.Stok).Scan(&p.ID)
	return p, err
}

func (r *productRepository) Update(id int, p models.Produk) (models.Produk, error) {
	_, err := r.db.Exec("UPDATE products SET nama = $1, harga = $2, stok = $3 WHERE id = $4",
		p.Nama, p.Harga, p.Stok, id)
	if err != nil {
		return p, err
	}
	p.ID = id
	return p, nil
}

func (r *productRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}