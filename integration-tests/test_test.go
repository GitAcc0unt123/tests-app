package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"tests_app/internal/handler"
	"tests_app/internal/models"
	"time"

	"github.com/lib/pq"
)

func (s *APITestSuite) TestTestCreate() {
	r := s.Require()

	test := models.Test{
		Title:                "test1",
		Description:          "description1",
		RandomQuestionsOrder: new(bool),
		QuestionsVisibility:  "ShowAll",
		Start:                time.Now(),
		End:                  time.Now().Add(5 * time.Minute),
		Duration_sec:         120,
	}

	testCreateData, err := json.Marshal(test)
	s.NoError(err)

	jwt, err := s.tokenManager.NewJWT(1)
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/test", bytes.NewBuffer(testCreateData))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var dbTest models.Test
	err = s.db.Get(&dbTest, `SELECT * FROM tests WHERE title = $1`, test.Title)
	s.NoError(err)

	r.Equal(1, dbTest.Id)
	r.Equal(test.Title, dbTest.Title)
	r.Equal(test.Description, dbTest.Description)
	r.Equal(test.RandomQuestionsOrder, dbTest.RandomQuestionsOrder)
	r.Equal(test.QuestionsVisibility, dbTest.QuestionsVisibility)
	//r.Equal(tests[0].Start, dbTest.Start)
	//r.Equal(tests[0].End, dbTest.End)
	r.Equal(test.Duration_sec, dbTest.Duration_sec)
}

func (s *APITestSuite) TestTestGetAll() {
	r := s.Require()

	tests := []models.Test{
		{
			Title:                "test1",
			Description:          "description1",
			RandomQuestionsOrder: new(bool),
			QuestionsVisibility:  "ShowAll",
			Start:                time.Now(),
			End:                  time.Now().Add(5 * time.Minute),
			Duration_sec:         120,
		},
		{
			Title:                "test2",
			Description:          "description2",
			RandomQuestionsOrder: new(bool),
			QuestionsVisibility:  "ShowAll",
			Start:                time.Now(),
			End:                  time.Now().Add(5 * time.Minute),
			Duration_sec:         120,
		},
	}

	_, err := s.db.Exec(`
	INSERT INTO tests (title, description, random_questions_order, questions_visibility, start_time, end_time, duration_sec)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		tests[0].Title,
		tests[0].Description,
		tests[0].RandomQuestionsOrder,
		tests[0].QuestionsVisibility,
		tests[0].Start,
		tests[0].End,
		tests[0].Duration_sec)
	s.NoError(err)

	_, err = s.db.Exec(`
	INSERT INTO tests (title, description, random_questions_order, questions_visibility, start_time, end_time, duration_sec)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		tests[1].Title,
		tests[1].Description,
		tests[1].RandomQuestionsOrder,
		tests[1].QuestionsVisibility,
		tests[1].Start,
		tests[1].End,
		tests[1].Duration_sec)
	s.NoError(err)

	request, _ := http.NewRequest("GET", "/api/test", nil)

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, request)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err := io.ReadAll(resp.Body)
	s.NoError(err)

	var responseTests []models.TestResponse
	err = json.Unmarshal(respData, &responseTests)
	s.NoError(err)

	r.Equal(tests[0].Title, responseTests[0].Title)
	r.Equal(tests[0].Description, responseTests[0].Description)
	r.Equal(tests[0].RandomQuestionsOrder, responseTests[0].RandomQuestionsOrder)
	r.Equal(tests[0].QuestionsVisibility, responseTests[0].QuestionsVisibility)
	//r.Equal(tests[0].Start, responseTests[0].Start)
	//r.Equal(tests[0].End, responseTests[0].End)
	r.Equal(tests[0].Duration_sec, responseTests[0].Duration_sec)

	r.Equal(tests[1].Title, responseTests[1].Title)
	r.Equal(tests[1].Description, responseTests[1].Description)
	r.Equal(tests[1].RandomQuestionsOrder, responseTests[1].RandomQuestionsOrder)
	r.Equal(tests[1].QuestionsVisibility, responseTests[1].QuestionsVisibility)
	//r.Equal(tests[1].Start, responseTests[1].Start)
	//r.Equal(tests[1].End, responseTests[1].End)
	r.Equal(tests[1].Duration_sec, responseTests[1].Duration_sec)
}

func (s *APITestSuite) TestTestGetFullById() {
	r := s.Require()

	user := models.User{
		Username: "username",
		Password: "qwerty123",
		Email:    "example@example.com",
	}

	test := models.Test{
		Title:                "test1",
		Description:          "description1",
		RandomQuestionsOrder: new(bool),
		QuestionsVisibility:  "ShowAll",
		Start:                time.Now(),
		End:                  time.Now().Add(5 * time.Minute),
		Duration_sec:         120,
	}

	questions := []models.Question{
		{
			TestId:      1,
			Text:        "question text",
			AnswerType:  "freeField",
			ShowAnswers: []string{},
			TrueAnswers: []string{"1999"},
		},
		{
			TestId:      1,
			Text:        "question text",
			AnswerType:  "oneSelect",
			ShowAnswers: []string{"1999", "2000", "2001"},
			TrueAnswers: []string{"1999"},
		},
	}

	_, err := s.db.Exec(`
	INSERT INTO tests (title, description, random_questions_order, questions_visibility, start_time, end_time, duration_sec)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		test.Title,
		test.Description,
		test.RandomQuestionsOrder,
		test.QuestionsVisibility,
		test.Start,
		test.End,
		test.Duration_sec)
	s.NoError(err)

	_, err = s.db.Exec(`
	INSERT INTO questions (test_id, text, answer_type, show_answers, true_answers)
	VALUES ($1, $2, $3, $4, $5)`,
		questions[0].TestId,
		questions[0].Text,
		questions[0].AnswerType,
		pq.Array(questions[0].ShowAnswers),
		pq.Array(questions[0].TrueAnswers))
	s.NoError(err)

	_, err = s.db.Exec(`
	INSERT INTO questions (test_id, text, answer_type, show_answers, true_answers)
	VALUES ($1, $2, $3, $4, $5)`,
		questions[1].TestId,
		questions[1].Text,
		questions[1].AnswerType,
		pq.Array(questions[1].ShowAnswers),
		pq.Array(questions[1].TrueAnswers))
	s.NoError(err)

	user.Password, err = s.hasher.Hash(user.Password)
	s.NoError(err)

	_, err = s.db.Exec(`
	INSERT INTO users (username, password, email, created_at)
	VALUES ($1, $2, $3, now())`,
		user.Username,
		user.Password,
		user.Email)
	s.NoError(err)

	jwt, err := s.tokenManager.NewJWT(1)
	s.NoError(err)

	request, _ := http.NewRequest("GET", "/api/test/1", nil)
	request.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, request)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err := io.ReadAll(resp.Body)
	s.NoError(err)

	var responseTests handler.GetTestByIdResponse
	err = json.Unmarshal(respData, &responseTests)
	s.NoError(err)
}
