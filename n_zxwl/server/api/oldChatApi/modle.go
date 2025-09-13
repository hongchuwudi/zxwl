package oldChatApi

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	SchoolID  int       `json:"school_id"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
}

func CreateTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS chat_messages (
        id INT AUTO_INCREMENT PRIMARY KEY,
        school_id INT NOT NULL,
        email VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        INDEX idx_school (school_id)
    )`
	_, err := db.Exec(query)
	return err
}
