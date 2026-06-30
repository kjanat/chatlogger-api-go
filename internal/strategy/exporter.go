// Package strategy implements the Strategy Pattern for the ChatLogger API.
// This file defines exporters that implement different strategies for data export
// formats (JSON, CSV). It follows the Strategy Pattern to allow for runtime selection
// of different export formats while maintaining a consistent interface.
package strategy

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/pkg/errors"
)

// Exporter defines the interface for export strategies.
// Implementations can export data in different formats (JSON, CSV, etc.).
type Exporter interface {
	// Export takes arbitrary data and exports it to a specific format.
	// It returns the exported data as bytes and any error encountered.
	Export(data interface{}) ([]byte, error)
}

// JSONExporter implements the Exporter interface for JSON format.
type JSONExporter struct{}

// Export exports data to JSON format.
func (j *JSONExporter) Export(data interface{}) ([]byte, error) {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal data to JSON")
	}
	return result, nil
}

// CSVExporter implements the Exporter interface for CSV format.
type CSVExporter struct{}

// Export exports data to CSV format.
// It expects 'data' to be a map with at least a 'chats' key containing an array of chats.
func (c *CSVExporter) Export(data interface{}) ([]byte, error) {
	chats, err := c.extractChatsFromData(data)
	if err != nil {
		return nil, err
	}

	writer, sb := c.createCSVWriter()
	if err := c.writeCSVHeader(writer); err != nil {
		return nil, err
	}

	if err := c.writeCSVData(writer, chats); err != nil {
		return nil, err
	}

	writer.Flush()
	return []byte(sb.String()), nil
}

// extractChatsFromData validates and extracts chats from the input data.
func (c *CSVExporter) extractChatsFromData(data interface{}) ([]domain.Chat, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data must be a map[string]interface{}")
	}

	chatsInterface, ok := dataMap["chats"]
	if !ok {
		return nil, fmt.Errorf("data must contain a 'chats' key")
	}

	chats, ok := chatsInterface.([]domain.Chat)
	if !ok {
		return nil, fmt.Errorf("'chats' must be an array of domain.Chat")
	}

	return chats, nil
}

// createCSVWriter creates a new CSV writer with string builder.
func (c *CSVExporter) createCSVWriter() (*csv.Writer, *strings.Builder) {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	return writer, &sb
}

// writeCSVHeader writes the CSV header row.
func (c *CSVExporter) writeCSVHeader(writer *csv.Writer) error {
	header := []string{
		"Chat ID", "Organization ID", "User ID", "Title", "Created At",
		"Message ID", "Role", "Content", "Timestamp", "Token Count", "Latency",
	}
	if err := writer.Write(header); err != nil {
		return errors.Wrap(err, "failed to write CSV header")
	}
	return nil
}

// writeCSVData writes all chat and message data to the CSV.
func (c *CSVExporter) writeCSVData(writer *csv.Writer, chats []domain.Chat) error {
	for _, chat := range chats {
		if err := c.writeChatToCSV(writer, chat); err != nil {
			return err
		}
	}
	return nil
}

// writeChatToCSV writes a single chat and its messages to the CSV.
func (c *CSVExporter) writeChatToCSV(writer *csv.Writer, chat domain.Chat) error {
	chatData := c.extractChatData(chat)

	if len(chat.Messages) == 0 {
		return c.writeEmptyChatRow(writer, chatData)
	}

	return c.writeChatWithMessages(writer, chat, chatData)
}

// extractChatData extracts common chat fields for CSV rows.
func (c *CSVExporter) extractChatData(chat domain.Chat) map[string]string {
	userID := "N/A"
	if chat.UserID != nil {
		userID = fmt.Sprintf("%d", *chat.UserID)
	}

	return map[string]string{
		"chatID":    fmt.Sprintf("%d", chat.ID),
		"orgID":     fmt.Sprintf("%d", chat.OrganizationID),
		"userID":    userID,
		"title":     chat.Title,
		"createdAt": chat.CreatedAt.Format(time.RFC3339),
	}
}

// writeEmptyChatRow writes a row for a chat with no messages.
func (c *CSVExporter) writeEmptyChatRow(writer *csv.Writer, chatData map[string]string) error {
	row := []string{
		chatData["chatID"], chatData["orgID"], chatData["userID"],
		chatData["title"], chatData["createdAt"], "", "", "", "", "", "",
	}
	if err := writer.Write(row); err != nil {
		return errors.Wrap(err, "failed to write CSV row for chat without messages")
	}
	return nil
}

// writeChatWithMessages writes rows for a chat with its messages.
func (c *CSVExporter) writeChatWithMessages(
	writer *csv.Writer, chat domain.Chat, chatData map[string]string,
) error {
	for _, message := range chat.Messages {
		messageData := c.extractMessageData(message)
		row := []string{
			chatData["chatID"], chatData["orgID"], chatData["userID"],
			chatData["title"], chatData["createdAt"],
			messageData["messageID"], messageData["role"], messageData["content"],
			messageData["timestamp"], messageData["tokenCount"], messageData["latency"],
		}
		if err := writer.Write(row); err != nil {
			return errors.Wrap(err, "failed to write CSV row for message")
		}
	}
	return nil
}

// extractMessageData extracts message fields for CSV rows.
func (c *CSVExporter) extractMessageData(message domain.Message) map[string]string {
	tokenCount := ""
	latency := ""

	if meta, err := message.GetMetadata(); err == nil && meta != nil {
		tokenCount = fmt.Sprintf("%d", meta.TokenCount)
		latency = fmt.Sprintf("%.2f", meta.ResponseTime)
	}

	return map[string]string{
		"messageID":  fmt.Sprintf("%d", message.ID),
		"role":       string(message.Role),
		"content":    message.Content,
		"timestamp":  message.CreatedAt.Format(time.RFC3339),
		"tokenCount": tokenCount,
		"latency":    latency,
	}
}
