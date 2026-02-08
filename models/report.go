package models

// BestSellingProduct represents the best selling product info
type BestSellingProduct struct {
	Nama      string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

// DailyReportResponse represents the daily sales report
type DailyReportResponse struct {
	TotalRevenue   int                `json:"total_revenue"`
	TotalTransaksi int                `json:"total_transaksi"`
	ProdukTerlaris BestSellingProduct `json:"produk_terlaris"`
}