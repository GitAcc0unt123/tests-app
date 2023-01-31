package models

type Question struct {
	Id          int      `json:"-"            db:"id"`
	TestId      int      `json:"test_id"      db:"test_id"      binding:"required"`
	Text        string   `json:"text"         db:"text"         binding:"required,min=1,max=255"`
	AnswerType  string   `json:"answer_type"  db:"answer_type"  binding:"required,oneof=freeField oneSelect manySelect"`
	ShowAnswers []string `json:"show_answers" db:"show_answers" binding:"required"`
	TrueAnswers []string `json:"true_answers" db:"true_answers" binding:"required"`
}

type GetQuestionResponse struct {
	Id          int      `json:"id"           db:"id"`
	TestId      int      `json:"test_id"      db:"test_id"`
	Text        string   `json:"text"         db:"text"`
	AnswerType  string   `json:"answer_type"  db:"answer_type"`
	ShowAnswers []string `json:"show_answers" db:"show_answers"`
}

type GetQuestionResponse2 struct {
	Id          int      `json:"id"           db:"id"`
	TestId      int      `json:"test_id"      db:"test_id"`
	Text        string   `json:"text"         db:"text"`
	AnswerType  string   `json:"answer_type"  db:"answer_type"`
	ShowAnswers []string `json:"show_answers" db:"show_answers"`
	TrueAnswers []string `json:"true_answers" db:"true_answers"`
}

type UpdateQuestionInput struct {
	Text        *string   `json:"text"`
	AnswerType  *string   `json:"answer_type"`
	ShowAnswers *[]string `json:"show_answers"`
	TrueAnswers *[]string `json:"true_answers"`
}
