package handlers

import (
	"category-api/services"
	"encoding/json"
	"net/http"
	"strings"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service}
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle /api/report/hari-ini
	if r.URL.Path == "/api/report/hari-ini" || strings.HasSuffix(r.URL.Path, "/api/report/hari-ini") {
		switch r.Method {
		case http.MethodGet:
			h.getDailyReport(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func (h *ReportHandler) getDailyReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}