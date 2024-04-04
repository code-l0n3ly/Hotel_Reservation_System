package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gorilla/mux"
)

type ReportHandler struct {
	db                *sql.DB
	ReportIdReference int64
	cache             map[string]Entities.Report // Cache to hold reports in memory
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{
		db:                db,
		ReportIdReference: 0,
		cache:             make(map[string]Entities.Report),
	}
}

func (ReportHandler *ReportHandler) GenerateUniqueReportID() string {
	ReportHandler.ReportIdReference++
	return fmt.Sprintf("%d", ReportHandler.ReportIdReference)
}

func (ReportHandler *ReportHandler) SetHighestReportID() {
	highestID := int64(0)
	for _, report := range ReportHandler.cache {
		reportID, err := strconv.ParseInt(report.ReportID, 10, 64)
		if err != nil {
			continue // Skip if the ReportID is not a valid integer
		}
		if reportID > highestID {
			highestID = reportID
		}
	}
	ReportHandler.ReportIdReference = highestID
}

func (ReportHandler *ReportHandler) LoadReports() error {
	rows, err := ReportHandler.db.Query(`SELECT ReportID, UserID, Type, CreateTime, Data FROM Report`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var report Entities.Report
		if err := rows.Scan(&report.ReportID, &report.UserID, &report.Type, &report.CreateTime, &report.Data); err != nil {
			return err
		}
		ReportHandler.cache[report.ReportID] = report
	}
	ReportHandler.SetHighestReportID()
	return rows.Err()
}

func (ReportHandler *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var report Entities.Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ReportHandler.LoadReports()
	err = report.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO Report (ReportID, UserID, Type, CreateTime, Data) VALUES (?, ?, ?, ?, ?)`
	_, err = ReportHandler.db.Exec(query, report.ReportID, report.UserID, report.Type, report.CreateTime, report.Data)
	if err != nil {
		http.Error(w, "Failed to create report", http.StatusInternalServerError)
		return
	}
	ReportHandler.LoadReports()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ReportHandler.cache[report.ReportID]) // Respond with the created report object
}

func (ReportHandler *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reportID := params["id"]

	var report Entities.Report
	query := `SELECT ReportID, UserID, Type, CreateTime, Data FROM Report WHERE ReportID = ?`
	err := ReportHandler.db.QueryRow(query, reportID).Scan(&report.ReportID, &report.UserID, &report.Type, &report.CreateTime, &report.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve report", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func (ReportHandler *ReportHandler) UpdateReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reportID := params["id"]
	ReportHandler.LoadReports()
	var report Entities.Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = report.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Report SET UserID = ?, Type = ?, CreateTime = ?, Data = ? WHERE ReportID = ?`
	_, err = ReportHandler.db.Exec(query, report.UserID, report.Type, report.CreateTime, report.Data, reportID)
	if err != nil {
		http.Error(w, "Failed to update report", http.StatusInternalServerError)
		return
	}
	ReportHandler.LoadReports()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Report updated successfully")
}

func (ReportHandler *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	reportID := params["id"]
	ReportHandler.LoadReports()
	query := `DELETE FROM Report WHERE ReportID = ?`
	_, err := ReportHandler.db.Exec(query, reportID)
	if err != nil {
		http.Error(w, "Failed to delete report", http.StatusInternalServerError)
		return
	}
	ReportHandler.LoadReports()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Report deleted successfully")
}
