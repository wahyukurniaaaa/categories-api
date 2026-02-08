package repositories

import (
	"category-api/models"
	"database/sql"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(items []models.CheckoutItem, products map[int]models.Produk) (models.Transaction, []models.TransactionDetail, error)
	GetTodayRevenue() (int, error)
	GetTodayTransactionCount() (int, error)
	GetTodayBestSellingProduct() (string, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) CreateTransaction(items []models.CheckoutItem, products map[int]models.Produk) (models.Transaction, []models.TransactionDetail, error) {
	// Start database transaction
	tx, err := r.db.Begin()
	if err != nil {
		return models.Transaction{}, nil, err
	}

	// Calculate total amount
	var totalAmount int
	for _, item := range items {
		product := products[item.ProductID]
		totalAmount += product.Harga * item.Quantity
	}

	// Insert transaction
	var transaction models.Transaction
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, total_amount, created_at",
		totalAmount,
	).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		tx.Rollback()
		return models.Transaction{}, nil, err
	}

	// Insert transaction details and update stock
	var details []models.TransactionDetail
	for _, item := range items {
		product := products[item.ProductID]
		subtotal := product.Harga * item.Quantity

		// Insert detail
		var detail models.TransactionDetail
		err = tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id, transaction_id, product_id, quantity, subtotal",
			transaction.ID, item.ProductID, item.Quantity, subtotal,
		).Scan(&detail.ID, &detail.TransactionID, &detail.ProductID, &detail.Quantity, &detail.Subtotal)
		if err != nil {
			tx.Rollback()
			return models.Transaction{}, nil, err
		}
		details = append(details, detail)

		// Update product stock
		_, err = tx.Exec(
			"UPDATE products SET stok = stok - $1 WHERE id = $2",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			tx.Rollback()
			return models.Transaction{}, nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.Transaction{}, nil, err
	}

	return transaction, details, nil
}

func (r *transactionRepository) GetTodayRevenue() (int, error) {
	var revenue sql.NullInt64
	today := time.Now().Format("2006-01-02")
	err := r.db.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE DATE(created_at) = $1",
		today,
	).Scan(&revenue)
	if err != nil {
		return 0, err
	}
	return int(revenue.Int64), nil
}

func (r *transactionRepository) GetTodayTransactionCount() (int, error) {
	var count int
	today := time.Now().Format("2006-01-02")
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM transactions WHERE DATE(created_at) = $1",
		today,
	).Scan(&count)
	return count, err
}

func (r *transactionRepository) GetTodayBestSellingProduct() (string, error) {
	var productName sql.NullString
	today := time.Now().Format("2006-01-02")
	err := r.db.QueryRow(`
		SELECT p.nama 
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = $1
		GROUP BY p.id, p.nama
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`, today).Scan(&productName)
	
	if err == sql.ErrNoRows {
		return "-", nil
	}
	if err != nil {
		return "", err
	}
	if !productName.Valid {
		return "-", nil
	}
	return productName.String, nil
}