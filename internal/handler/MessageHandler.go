package Handlers

import (
	"database/sql"
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
	rows, err := handler.db.Query(`SELECT * FROM Chat`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var createTime []byte
		var Chat Entities.Chat
		if err := rows.Scan(&Chat.ChatID, &Chat.SenderID, &Chat.ReceiverID, &createTime); err != nil {
			return err
		}
		Chat.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
		rows, err := handler.db.Query(`SELECT MessageID, ChatID, SenderID, Content, CreateTime FROM Message WHERE ChatID = ?`, Chat.ChatID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var createTime []byte
			var message Entities.Message
			if err := rows.Scan(&message.MessageID, &message.ChatID, &message.SenderID, &message.Content, &createTime); err != nil {
				return err
			}
			message.CreateTime, _ = time.Parse("2006-01-02 15:04:05", string(createTime))
			Chat.Messages = append(Chat.Messages, message)
		}
		handler.cache[Chat.ChatID] = Chat
	}
	return rows.Err()

}

// function to check if there is a chat between two users, if there is return the chatID otherwise return an empty string
func (handler *MessageHandler) GetChatID(senderID string, receiverID string) string {
	handler.LoadMessages()
	for _, chat := range handler.cache {
		if (chat.SenderID == senderID && chat.ReceiverID == receiverID) || (chat.SenderID == receiverID && chat.ReceiverID == senderID) {
			return chat.ChatID
		}
	}
	//Create a new chat
	query := `INSERT INTO Chat (SenderID, ReceiverID) VALUES (?, ?)`
	_, err := handler.db.Exec(query, senderID, receiverID)
	if err != nil {
		return ""
	}
	handler.LoadMessages()
	for _, chat := range handler.cache {
		if (chat.SenderID == senderID && chat.ReceiverID == receiverID) || (chat.SenderID == receiverID && chat.ReceiverID == senderID) {
			return chat.ChatID
		}
	}
	return ""
}
func (handler *MessageHandler) SendMessage(c *gin.Context) {
	var message Entities.Message
	err := c.BindJSON(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	handler.LoadMessages()
	message.ChatID = handler.GetChatID(message.SenderID, message.ReceiverID)
	if message.ChatID == "" {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to create chat",
		})
		return
	}
	query := `INSERT INTO Message (ChatID, SenderID, Content) VALUES (?, ?, ?)`
	_, err = handler.db.Exec(query, message.ChatID, message.SenderID, message.Content)
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

func (handler *MessageHandler) GetChat(c *gin.Context) {
	var chatRequest struct {
		SenderID   string `json:"senderID"`
		ReceiverID string `json:"receiverID"`
	}

	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	chatID := ""
	for _, chat := range handler.cache {
		if (chat.SenderID == chatRequest.SenderID && chat.ReceiverID == chatRequest.ReceiverID) || (chat.SenderID == chatRequest.ReceiverID && chat.ReceiverID == chatRequest.SenderID) {
			chatID = chat.ChatID
		}
	}
	if chatID == "" {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Chat not found",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Chat retrieved successfully",
		Data:    handler.cache[chatID],
	})
}

// Get chat by Chat ID
func (handler *MessageHandler) GetChatByID(c *gin.Context) {
	chatID := c.Param("id")
	handler.LoadMessages()
	chat, exists := handler.cache[chatID]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Chat not found",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Chat retrieved successfully",
		Data:    chat,
	})
}

// Get chat by sender ID
func (handler *MessageHandler) GetChatBySenderID(c *gin.Context) {
	senderID := c.Param("id")
	handler.LoadMessages()
	var chats []Entities.Chat
	for _, chat := range handler.cache {
		if chat.SenderID == senderID {
			chats = append(chats, chat)
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Chats retrieved successfully",
		Data:    chats,
	})
}
