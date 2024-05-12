package Handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	Entities "GraduationProject.com/m/internal/model"
	"github.com/gin-gonic/gin"
)

type FinancialTransactionHandler struct {
	db    *sql.DB
	cache map[string]Entities.FinancialTransaction // Cache to hold transactions in memory
}

func NewFinancialTransactionHandler(db *sql.DB) *FinancialTransactionHandler {
	return &FinancialTransactionHandler{
		db:    db,
		cache: make(map[string]Entities.FinancialTransaction),
	}
}

func (handler *FinancialTransactionHandler) LoadTransactions() error {
	handler.cache = make(map[string]Entities.FinancialTransaction)
	rows, err := handler.db.Query(`SELECT TransactionID, UserID, BookingID, PaymentMethod, Amount, CreateTime FROM FinancialTransaction`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Entities.FinancialTransaction
		var createTime []byte
		if err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.BookingID, &transaction.PaymentMethod, &transaction.Amount, &createTime); err != nil {
			fmt.Println(err.Error())
			return err
		}
		transaction.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		handler.cache[transaction.TransactionID] = transaction
	}
	return rows.Err()
}

func (handler *FinancialTransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction Entities.FinancialTransaction
	handler.LoadTransactions()

	err := c.BindJSON(&transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	query := `INSERT INTO FinancialTransaction (UserID, BookingID, PaymentMethod, Amount) VALUES (?, ?, ?, ?)`
	result, err := handler.db.Exec(query, transaction.UserID, transaction.BookingID, transaction.PaymentMethod, transaction.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create transaction" + err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	transaction.TransactionID = strconv.FormatInt(id, 10)
	handler.LoadTransactions()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Transaction created successfully", "data": handler.cache[transaction.TransactionID]})
}

func (handler *FinancialTransactionHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	handler.LoadTransactions()

	transaction, exists := handler.cache[transactionID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transaction retrieved successfully", "data": transaction})
}

func (handler *FinancialTransactionHandler) UpdateTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	handler.LoadTransactions()

	var newInfoTransaction Entities.FinancialTransaction
	oldInfoTransaction := handler.cache[transactionID]

	err := c.BindJSON(&newInfoTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if newInfoTransaction.PaymentMethod != "" {
		oldInfoTransaction.PaymentMethod = newInfoTransaction.PaymentMethod
	}
	if newInfoTransaction.Amount != 0 {
		oldInfoTransaction.Amount = newInfoTransaction.Amount
	}

	query := `UPDATE FinancialTransaction SET PaymentMethod = ?, Amount = ? WHERE TransactionID = ?`
	_, err = handler.db.Exec(query, oldInfoTransaction.PaymentMethod, oldInfoTransaction.Amount, oldInfoTransaction.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to update transaction" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transaction updated successfully", "Data": oldInfoTransaction})
}

func (handler *FinancialTransactionHandler) DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	handler.LoadTransactions()

	query := `DELETE FROM FinancialTransaction WHERE TransactionID = ?`
	_, err := handler.db.Exec(query, transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to delete transaction" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transaction deleted successfully", "data": handler.cache[transactionID]})
	handler.LoadTransactions()
}

// GET all the transactions for a user
func (handler *FinancialTransactionHandler) GetTransactionsByUserID(c *gin.Context) {
	userID := c.Param("id")
	handler.LoadTransactions()
	var userTransactions []Entities.FinancialTransaction
	for _, transaction := range handler.cache {
		if transaction.UserID == userID {
			userTransactions = append(userTransactions, transaction)
		}
	}
	if len(userTransactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No transactions found for this user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transactions retrieved successfully", "data": userTransactions})
}
