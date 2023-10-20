package repositories

import (
	"context"
	"log"
	"my_chatbot/infrastructures"
	"my_chatbot/models"
	"time"
)

type MessageRepository struct{}

func (m *MessageRepository) StoreMessage(ctx context.Context, userid int, usermessage, response string, timestamp time.Time, path string) {
	db, err := infrastructures.GetDB()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO messages (userid, usermessage, response, timestamp, path) VALUES (?, ?, ?, ?, ?)", userid, usermessage, response, timestamp.Format("2006-01-02 15:04:05"), path)
	if err != nil {
		panic(err)
	}
}

func (m *MessageRepository) GetMessage(ctx context.Context, userid int, limit int, offset int) ([]models.Message, error) {
	if limit == 0 {
		limit = 10
	}
	db, _ := infrastructures.GetDB()
	rows, err := db.Query("SELECT userid, usermessage, response, path, timestamp FROM messages WHERE userid = ? order by timestamp desc limit ? offset ?", userid, limit, offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.UserID, &msg.UserMessage, &msg.Response, &msg.Path, &msg.Timestamp)
		if err != nil {
			panic(err)
		}
		messages = append(messages, msg)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return messages, err
}
