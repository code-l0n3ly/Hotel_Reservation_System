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

type MessageHandler struct {
	db                 *sql.DB
	MessageIdReference int64
	cache              map[string]Entities.Message // Cache to hold messages in memory
}

func NewMessageHandler(db *sql.DB) *MessageHandler {
	return &MessageHandler{
		db:                 db,
		MessageIdReference: 0,
		cache:              make(map[string]Entities.Message),
	}
}

func (handler *MessageHandler) GenerateUniqueMessageID() string {
	handler.MessageIdReference++
	return fmt.Sprintf("%d", handler.MessageIdReference)
}

func (handler *MessageHandler) SetHighestMessageID() {
	highestID := int64(0)
	for _, message := range handler.cache {
		messageID, err := strconv.ParseInt(message.MessageID, 10, 64)
		if err != nil {
			continue // Skip if the MessageID is not a valid integer
		}
		if messageID > highestID {
			highestID = messageID
		}
	}
	handler.MessageIdReference = highestID
}

func (handler *MessageHandler) LoadMessages() error {
	rows, err := handler.db.Query(`SELECT MessageID, Content, CreateTime, ReceiverID, SenderID FROM Message`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var message Entities.Message
		if err := rows.Scan(&message.MessageID, &message.Content, &message.CreateTime, &message.ReceiverID, &message.SenderID); err != nil {
			return err
		}
		handler.cache[message.MessageID] = message
	}
	handler.SetHighestMessageID()
	return rows.Err()
}

func (handler *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message Entities.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handler.LoadMessages()

	query := `INSERT INTO Message (MessageID, Content, CreateTime, ReceiverID, SenderID) VALUES (?, ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, message.MessageID, message.Content, message.CreateTime, message.ReceiverID, message.SenderID)
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	handler.LoadMessages()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(handler.cache[message.MessageID]) // Respond with the created message object
}

func (handler *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messageID := params["id"]

	var message Entities.Message
	query := `SELECT MessageID, Content, CreateTime, ReceiverID, SenderID FROM Message WHERE MessageID = ?`
	err := handler.db.QueryRow(query, messageID).Scan(&message.MessageID, &message.Content, &message.CreateTime, &message.ReceiverID, &message.SenderID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to retrieve message", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(message)
}

func (handler *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messageID := params["id"]
	handler.LoadMessages()
	var message Entities.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE Message SET Content = ?, CreateTime = ?, ReceiverID = ?, SenderID = ? WHERE MessageID = ?`
	_, err = handler.db.Exec(query, message.Content, message.CreateTime, message.ReceiverID, message.SenderID, messageID)
	if err != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}
	handler.LoadMessages()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Message updated successfully")
}

func (handler *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messageID := params["id"]
	handler.LoadMessages()
	query := `DELETE FROM Message WHERE MessageID = ?`
	_, err := handler.db.Exec(query, messageID)
	if err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}
	handler.LoadMessages()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Message deleted successfully")
}
