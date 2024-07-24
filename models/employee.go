package models

import (
	"time"
)

type Employee struct {
	Name      string
	Category  string
	JoinDate  time.Time
	BirthDate time.Time
	Email     string
}

type LogEntry struct {
	ID      int    `db:"id"`
	Type    string `db:"type"`
	Message string `db:"message"`
	Date    int64  `db:"date"`
}

type UnsentEmail struct {
	ID     int    `db:"id"`
	Email  string `db:"email"`
	Reason string `db:"reason"`
}
