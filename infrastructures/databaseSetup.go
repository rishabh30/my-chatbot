package infrastructures

import (
	"database/sql"
	"fmt"
	"sync"
)

var (
	dbInstance        *sql.DB
	dbConnectionError error
	dbsOnce           sync.Once
)

func InitializeDatabase() error {
	dbsOnce.Do(func() {
		dbInstance, dbConnectionError = sql.Open("sqlite3", "./storage/messages.db")
		if dbConnectionError != nil {
			return
		}

		_, dbConnectionError = dbInstance.Exec(`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY, 
			userid INTEGER, 
			usermessage TEXT, 
			response TEXT, 
			timestamp TEXT,
			path TEXT
		)`)
	})
	return dbConnectionError
}

func GetDB() (*sql.DB, error) {
	if dbInstance == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return dbInstance, nil
}
