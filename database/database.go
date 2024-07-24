package database

import (
	"log"
	"mailsys/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func Init() {
	var err error
	DB, err = sqlx.Connect("sqlite3", "./mailsys.db")
	if err != nil {
		log.Fatalln(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS log_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		message TEXT,
		date TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS unsent_emails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT,
		reason TEXT
	);
	`
	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
}

func InsertLogEntry(logEntry models.LogEntry) error {
	query := `INSERT INTO log_entries (type, message, date) VALUES (:type, :message, :date)`
	_, err := DB.NamedExec(query, logEntry)
	return err
}

func InsertUnsentEmail(unsentEmail models.UnsentEmail) error {
	query := `INSERT INTO unsent_emails (email, reason) VALUES (:email, :reason)`
	_, err := DB.NamedExec(query, unsentEmail)
	return err
}

func GetUnsentEmails() ([]models.UnsentEmail, error) {
	var unsentEmails []models.UnsentEmail
	err := DB.Select(&unsentEmails, "SELECT * FROM unsent_emails")
	return unsentEmails, err
}

func DeleteUnsentEmail(id int) error {
	query := `DELETE FROM unsent_emails WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}
