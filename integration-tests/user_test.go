package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"tests_app/internal/models"
	"time"
)

func (s *APITestSuite) TestUserSignUp() {
	r := s.Require()

	username, email, password := "TestUser", "test@test.com", "qwerty123"
	signUpData := fmt.Sprintf(`{"username":"%s","email":"%s","password":"%s"}`, username, email, password)

	req, _ := http.NewRequest("POST", "/api/auth/sign-up", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var user models.User
	err := s.db.Get(&user, `SELECT * FROM users WHERE username = $1`, username)
	s.NoError(err)

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	r.Equal("", user.Name)
	r.Equal(username, user.Username)
	r.Equal(passwordHash, user.Password)
	r.Equal(email, user.Email)

	var activated_at *time.Time = nil
	r.Equal(activated_at, user.ActivatedAt)
}
