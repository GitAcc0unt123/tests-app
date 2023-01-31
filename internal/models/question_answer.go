package models

import (
	"time"
)

type QuestionAnswer struct {
	UserId     int       `json:"-"           db:"user_id"`
	QuestionId int       `json:"question_id" db:"question_id" binding:"required"`
	Answer     []string  `json:"answer"      db:"answer"      binding:"required"`
	Time       time.Time `json:"time"        db:"time"`
}

type UpsertQuestionAnswerInput struct {
	QuestionId int      `json:"question_id" binding:"required"`
	Answer     []string `json:"answer"      binding:"required"`
}
