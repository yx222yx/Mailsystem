package models

import "time"

type Employee struct {
	Name      string
	Category  string
	JoinDate  time.Time
	BirthDate time.Time
	Email     string
}
