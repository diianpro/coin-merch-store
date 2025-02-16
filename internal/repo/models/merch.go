package models

import "time"

type Merch struct {
	Name   string    `db:"name"`
	Amount int64     `db:"amount"`
	Date   time.Time `db:"date"`
}
