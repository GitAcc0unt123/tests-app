package models

import (
	"time"
)

/*type PgDuration struct {
	time.Duration `json:""`
}
func (d *PgDuration) MarshalJSON() ([]byte, error) {}
func (d *PgDuration) UnmarshalJSON(data []byte) error {}
func (d PgDuration) Value() (driver.Value, error) {}
func (d *PgDuration) Scan(raw interface{}) error {}*/

type Test struct {
	Id                   int       `json:"-"                      db:"id"`
	Title                string    `json:"title"                  db:"title"                  binding:"required,min=1,max=255"`
	Description          string    `json:"description"            db:"description"            binding:"max=2048"`
	RandomQuestionsOrder *bool     `json:"random_questions_order" db:"random_questions_order" binding:"required"`
	QuestionsVisibility  string    `json:"questions_visibility"   db:"questions_visibility"   binding:"required,oneof=ShowOneByOne ShowAll"`
	Start                time.Time `json:"start_time"             db:"start_time"             binding:"required,ltfield=End"` // time_format:"yyyy-mm-dd hh:mm:ss"`
	End                  time.Time `json:"end_time"               db:"end_time"               binding:"required"`             // time_format:"yyyy-mm-dd hh:mm:ss"`
	Duration_sec         int       `json:"duration_sec"           db:"duration_sec"           binding:"required" swaggertype:"integer"`
}

type TestResponse struct {
	Id                   int       `json:"id"                     db:"id"`
	Title                string    `json:"title"                  db:"title"`
	Description          string    `json:"description"            db:"description"`
	RandomQuestionsOrder *bool     `json:"random_questions_order" db:"random_questions_order"`
	QuestionsVisibility  string    `json:"questions_visibility"   db:"questions_visibility"`
	Start                time.Time `json:"start_time"             db:"start_time"`
	End                  time.Time `json:"end_time"               db:"end_time"`
	Duration_sec         int       `json:"duration_sec"           db:"duration_sec"`
}

type UpdateTestInput struct {
	Title                *string    `json:"title"                 binding:"min=1,max=255"`
	Description          *string    `json:"description"           binding:"max=2048"`
	RandomQuestionsOrder *bool      `json:"random_questions_order"`
	QuestionsVisibility  *string    `json:"questions_visibility"  binding:"oneof=ShowOneByOne ShowAll"`
	Start                *time.Time `json:"start_time"`
	End                  *time.Time `json:"end_time"`
	Duration_sec         *int       `json:"duration_sec"          swaggertype:"integer"`
}
