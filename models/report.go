package models

// DailyReportResponse represents the daily sales report
type DailyReportResponse struct {
	TotalRevenue   int    `json:"total_revenue"`
	TotalTransaksi int    `json:"total_transaksi"`
	ProdukTerlaris string `json:"produk_terlaris"`
}