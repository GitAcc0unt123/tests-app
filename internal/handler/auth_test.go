package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"tests_app/internal/config"
	"tests_app/internal/models"
	"tests_app/internal/service"

	service_mocks "tests_app/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockUser, user models.User)
	mockStub := func(r *service_mocks.MockUser, user models.User) {}

	tests := []struct {
		name                 string
		requestBody          string
		serviceMockParam     models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			requestBody: `{"username": "username", "password": "qwerty1234", "name": "Test Name", "email": "example@example.com"}`,
			serviceMockParam: models.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty1234",
				Email:    "example@example.com",
			},
			mockBehavior: func(r *service_mocks.MockUser, user models.User) {
				r.EXPECT().Create(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:        "Ok",
			requestBody: `{"username": "username", "password": "qwerty1234", "email": "example@example.com"}`,
			serviceMockParam: models.User{
				Username: "username",
				Password: "qwerty1234",
				Email:    "example@example.com",
			},
			mockBehavior: func(r *service_mocks.MockUser, user models.User) {
				r.EXPECT().Create(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1}`,
		},

		{
			name:                 "Empty Input",
			requestBody:          ``,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{`,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:                 "Invalid JSON",
			requestBody:          `{username": "username"}`,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},

		{
			name:                 "Empty Input",
			requestBody:          `{}`,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:                 "required fields",
			requestBody:          `{"password": "qwerty1234"}`,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:                 "required fields",
			requestBody:          `{"username": "username", "password": "qwerty1234"}`,
			serviceMockParam:     models.User{},
			mockBehavior:         mockStub,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},

		{
			name:        "Service Error. предположим unique constraint",
			requestBody: `{"username": "username", "password": "qwerty1234", "email": "example@example.com"}`,
			serviceMockParam: models.User{
				Username: "username",
				Password: "qwerty1234",
				Email:    "example@example.com",
			},
			mockBehavior: func(r *service_mocks.MockUser, user models.User) {
				r.EXPECT().Create(user).Return(0, errors.New(""))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockUser(c)
			test.mockBehavior(repo, test.serviceMockParam)

			services := &service.Service{User: repo}
			handler := Handler{services, config.Config{}} // ???

			// Init Endpoint
			r := gin.New()
			r.POST("api/auth/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/auth/sign-up",
				bytes.NewBufferString(test.requestBody),
			)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
