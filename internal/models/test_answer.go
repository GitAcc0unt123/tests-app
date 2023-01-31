package models

import "time"

type TestAnswer struct {
	UserId       int       `json:"-"       db:"user_id"`
	TestId       int       `json:"test_id" db:"test_id"       binding:"required"`
	CompleteTime time.Time `json:"-"       db:"complete_time"`
}
