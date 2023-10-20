package interfaces

import (
	"context"
	"my_chatbot/models"
	"time"
)

type IMessageRepository interface {
	GetMessage(ctx context.Context, userid int, limit int, offset int) ([]models.Message, error)
	StoreMessage(ctx context.Context, userid int, usermessage, response string, timestamp time.Time, path string)
}
