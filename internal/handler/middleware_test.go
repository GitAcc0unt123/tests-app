package handler

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_getUserId(t *testing.T) {
	testTable := []struct {
		name          string
		context       func() *gin.Context
		expectedValue int
		expectedError error
	}{
		{
			name: "id undefined",
			context: func() *gin.Context {
				return &gin.Context{}
			},
			expectedValue: 0,
			expectedError: errors.New("user id not found"),
		},
		{
			name: "invalid type",
			context: func() *gin.Context {
				context := &gin.Context{}
				context.Set("userId", "1")
				return context
			},
			expectedValue: 0,
			expectedError: errors.New("user id is of invalid type"),
		},
		{
			name: "ok",
			context: func() *gin.Context {
				context := &gin.Context{}
				context.Set("userId", 1)
				return context
			},
			expectedValue: 1,
			expectedError: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			userId, err := getUserId(test.context())

			assert.Equal(t, test.expectedValue, userId)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
