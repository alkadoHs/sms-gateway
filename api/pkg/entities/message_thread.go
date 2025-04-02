package entities

import (
	"time"

	"github.com/google/uuid"
)

// MessageThread represents a message thread between 2 phone numbers
type MessageThread struct {
	ID                 uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid;" example:"32343a19-da5e-4b1b-a767-3298a73703ca"`
	Owner              string        `json:"owner" example:"+18005550199"`
	Contact            string        `json:"contact" example:"+18005550100"`
	IsArchived         bool          `json:"is_archived" example:"false"`
	UserID             UserID        `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	Color              string        `json:"color" example:"indigo"`
	Status             MessageStatus `json:"status" example:"PENDING"`
	LastMessageContent *string       `json:"last_message_content" example:"This is a sample message content"`
	LastMessageID      *uuid.UUID    `json:"last_message_id" example:"32343a19-da5e-4b1b-a767-3298a73703ca"`
	CreatedAt          time.Time     `json:"created_at" example:"2022-06-05T14:26:09.527976+03:00"`
	UpdatedAt          time.Time     `json:"updated_at" example:"2022-06-05T14:26:09.527976+03:00"`
	OrderTimestamp     time.Time     `json:"order_timestamp" example:"2022-06-05T14:26:09.527976+03:00"`
}

// Update a message thread after a message event
func (thread *MessageThread) Update(timestamp time.Time, messageID uuid.UUID, content string, status MessageStatus) *MessageThread {
	thread.OrderTimestamp = timestamp
	thread.LastMessageID = &messageID
	thread.Status = status
	thread.LastMessageContent = &content
	return thread
}

// UpdateArchive sets a message thread as archived
func (thread *MessageThread) UpdateArchive(isArchived bool) *MessageThread {
	thread.IsArchived = isArchived
	return thread
}

// HasLastMessage checks the last message in a thread by ID
func (thread *MessageThread) HasLastMessage(id uuid.UUID) bool {
	if thread.LastMessageID == nil {
		return false
	}
	return *thread.LastMessageID == id
}
