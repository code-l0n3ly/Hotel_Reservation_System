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

type FinancialTransactionHandler struct {
	db                     *sql.DB
	TransactionIdReference int64
	cache                  map[string]Entities.FinancialTransaction // Cache to hold transactions in memory
}

func NewFinancialTransactionHandler(db *sql.DB) *FinancialTransactionHandler {
	return &FinancialTransactionHandler{
		db:                     db,
		TransactionIdReference: 0,
		cache:                  make(map[string]Entities.FinancialTransaction),
	}
}

func (handler *FinancialTransactionHandler) GenerateUniqueTransactionID() string {
	handler.TransactionIdReference++
	return fmt.Sprintf("%d", handler.TransactionIdReference)
}

func (handler *FinancialTransactionHandler) SetHighestTransactionID() {
	highestID := int64(0)
	for _, transaction := range handler.cache {
		transactionID, err := strconv.ParseInt(transaction.TransactionID, 10, 64)
		if err != nil {
			continue // Skip if the TransactionID is not a valid integer
		}
		if transactionID > highestID {
			highestID = transactionID
		}
	}
	handler.TransactionIdReference = highestID
}

func (handler *FinancialTransactionHandler) LoadTransactions() error {
	rows, err := handler.db.Query(`SELECT TransactionID, UserID, UnitID, PaymentMethod, Amount, CreateTime FROM FinancialTransaction`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Entities.FinancialTransaction
		if err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.UnitID, &transaction.PaymentMethod, &transaction.Amount, &transaction.CreateTime); err != nil {
			return err
		}
		handler.cache[transaction.TransactionID] = transaction
	}
	handler.SetHighestTransactionID()
	return rows.Err()
}

func (handler *FinancialTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Entities.FinancialTransaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handler.LoadTransactions()

	query := `INSERT INTO FinancialTransaction (TransactionID, UserID, UnitID, PaymentMethod, Amount, CreateTime) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, transaction.TransactionID, transaction.UserID, transaction.UnitID, transaction.PaymentMethod, transaction.Amount, transaction.CreateTime)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}
	handler.LoadTransactions()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.cache[transaction.TransactionID]) // Respond with the created transaction object
}

func (handler *FinancialTransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionID := params["id"]

	var transaction Entities.FinancialTransaction
	query := `SELECT TransactionID, UserID, UnitID, PaymentMethod, Amount, CreateTime FROM FinancialTransaction WHERE TransactionID = ?`
	err := handler.db.QueryRow(query, transactionID).Scan(&transaction.TransactionID, &transaction.UserID, &transaction.UnitID, &transaction.PaymentMethod, &transaction.Amount, &transaction.CreateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve transaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (handler *FinancialTransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionID := params["id"]
	handler.LoadTransactions()
	var transaction Entities.FinancialTransaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE FinancialTransaction SET UserID = ?, UnitID = ?, PaymentMethod = ?, Amount = ?, CreateTime = ? WHERE TransactionID = ?`
	_, err = handler.db.Exec(query, transaction.UserID, transaction.UnitID, transaction.PaymentMethod, transaction.Amount, transaction.CreateTime, transactionID)
	if err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}
	handler.LoadTransactions()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Transaction updated successfully")
}

func (handler *FinancialTransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionID := params["id"]
	handler.LoadTransactions()
	query := `DELETE FROM FinancialTransaction WHERE TransactionID = ?`
	_, err := handler.db.Exec(query, transactionID)
	if err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}
	handler.LoadTransactions()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Transaction deleted successfully")
}
