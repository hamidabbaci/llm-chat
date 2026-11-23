package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=" + dbPassword + " dbname=llmchat sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("خطا در اتصال به دیتابیس:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("دیتابیس در دسترس نیست:", err)
	}

	log.Println("اتصال به دیتابیس برقرار شد")
}
func saveMessage(role, content string) error {
	_, err := db.Exec(
		"INSERT INTO messages (role, content) VALUES ($1, $2)",
		role, content,
	)
	return err
}

func loadHistory() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT role, content FROM messages ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var role, content string
		if err := rows.Scan(&role, &content); err != nil {
			return nil, err
		}
		history = append(history, map[string]interface{}{
			"role":  role,
			"parts": []map[string]string{{"text": content}},
		})
	}

	return history, nil
}
