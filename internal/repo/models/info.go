package models

import "time"

type Info struct {
	FromUserID int64     `db:"from_user_id"`
	ToUserID   int64     `db:"to_user_id"`
	Amount     int64     `db:"amount"`
	Date       time.Time `db:"date"`
}
