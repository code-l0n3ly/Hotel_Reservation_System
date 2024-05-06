package Handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	Entities "GraduationProject.com/m/internal/model"
)

type MessageHandler struct {
	db    *sql.DB
	cache map[string]Entities.Chat // Cache to hold messages in memory
}

func NewMessageHandler(db *sql.DB) *MessageHandler {
	return &MessageHandler{
		db:    db,
		cache: make(map[string]Entities.Chat),
	}
}

func (handler *MessageHandler) LoadMessages() error {
	rows, err := handler.db.Query(`SELECT ChatID, CreateTime FROM Chat`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var createTime []byte
		var Chat Entities.Chat
		if err := rows.Scan(&Chat.ChatID, &createTime); err != nil {
			return err
		}
		Chat.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		rows, err := handler.db.Query(`SELECT MessageID, SenderID, Content, CreateTime FROM Message WHERE ChatID = ?`, Chat.ChatID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var createTime []byte
			var message Entities.Message
			if err := rows.Scan(&message.MessageID, &message.Content, &createTime, &message.SenderID); err != nil {
				return err
			}
			message.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		}
		handler.cache[Chat.ChatID] = Chat
	}
	return rows.Err()

}

func (handler *MessageHandler) CreateMessage(c *gin.Context) {
	var message Entities.Message
	err := json.NewDecoder(c.Request.Body).Decode(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	handler.LoadMessages()

	query := `INSERT INTO Message (Content, CreateTime, ReceiverID, SenderID) VALUES ( ?, ?, ?, ?)`
	_, err = handler.db.Exec(query, message.Content, message.CreateTime, message.ReceiverID, message.SenderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to create message",
		})
		return
	}
	handler.LoadMessages()
	c.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "Message created successfully",
		Data:    handler.cache[message.MessageID],
	})
}

func (handler *MessageHandler) GetMessage(c *gin.Context) {
	messageID := c.Param("id")

	message, exists := handler.cache[messageID]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Message not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message retrieved successfully",
		Data:    message,
	})
}

func (handler *MessageHandler) UpdateMessage(c *gin.Context) {
	messageID := c.Param("id")
	handler.LoadMessages()

	var message Entities.Message
	err := json.NewDecoder(c.Request.Body).Decode(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Start building the query
	query := "UPDATE Message SET "
	args := []interface{}{}

	if message.Content != "" {
		query += "Content = ?, "
		args = append(args, message.Content)
	}

	if !message.CreateTime.IsZero() {
		query += "CreateTime = ?, "
		args = append(args, message.CreateTime)
	}

	if message.ReceiverID != "" {
		query += "ReceiverID = ?, "
		args = append(args, message.ReceiverID)
	}

	if message.SenderID != "" {
		query += "SenderID = ?, "
		args = append(args, message.SenderID)
	}

	// Remove the last comma and space, and add the WHERE clause
	query = query[:len(query)-2] + " WHERE MessageID = ?"
	args = append(args, messageID)

	_, err = handler.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to update message",
		})
		return
	}

	handler.LoadMessages()
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message updated successfully",
	})
}

func (handler *MessageHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")
	handler.LoadMessages()
	query := `DELETE FROM Message WHERE MessageID = ?`
	_, err := handler.db.Exec(query, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to delete message",
		})
		return
	}
	handler.LoadMessages()
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message deleted successfully",
	})
}
