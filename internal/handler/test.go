package handler

import (
	"net/http"
	"strconv"
	"tests_app/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// createTest godoc
//
// @Summary     Create a test
// @Security    ApiKeyAuth
// @Description Create a test
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       data body     models.Test true "test info"
// @Success     201  {object} idResponse
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /test [post]
func (h *Handler) createTest(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.Test
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	testId, err := h.services.Test.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusCreated, idResponse{
		Id: testId,
	})
}

// getAllTests godoc
//
// @Summary     Get all tests
// @Description Get all tests
// @Tags        test
// @Accept      json
// @Produce     json
// @Success     200 {object} []models.TestResponse
// @Failure     500 {object} errorResponse
// @Router      /test [get]
func (h *Handler) getAllTests(c *gin.Context) {
	tests, err := h.services.Test.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, tests)
}

// getTestById godoc
//
// @Summary     Get test by id
// @Description Get test by id
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       id  path     int           true "A unique integer value identifying this test."
// @Success     200 {object} models.Test   "test info"
// @Failure     404 {object} errorResponse "incorrect id"
// @Failure     403 {object} errorResponse "forbidden"
// @Failure     500 {object} errorResponse "todo"
// @Router      /test/{id} [get]
func (h *Handler) getTestById(c *gin.Context) {
	testId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "incorrect id: "+c.Param("id"), nil)
		return
	}

	test, err := h.services.Test.GetById(testId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, test)
}

// getFullTestById godoc
//
// @Summary     Get test by id
// @Security    ApiKeyAuth
// @Description Тест с вопросами и ответами если есть
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       id  path     int                  true "A unique integer value identifying this test."
// @Success     200 {object} GetTestByIdResponse  "test info"
// @Success     200 {object} GetTestByIdResponse2 "test info"
// @Failure     404 {object} errorResponse        "incorrect id"
// @Failure     403 {object} errorResponse        "forbidden"
// @Failure     500 {object} errorResponse        "todo"
// @Router      /test/{id}/full [get]
func (h *Handler) getFullTestById(c *gin.Context) {
	testId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "incorrect id: "+c.Param("id"), nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error(), err)
		return
	}

	test, err := h.services.Test.GetById(testId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	answers, err := h.services.QuestionAnswer.GetAnswersByTestId(userId, testId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	testAnswer, err := h.services.TestAnswer.Get(userId, testId)
	if err != nil {
		questions, err := h.services.Question.GetAll(userId, testId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
			return
		}

		c.JSON(http.StatusOK, GetTestByIdResponse{
			Test:      test,
			Questions: questions,
			Answers:   answers,
		})
		return
	}

	questions, err := h.services.Question.GetAllWithAnswer(userId, testId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, GetTestByIdResponse2{
		Test:       test,
		TestAnswer: testAnswer,
		Questions:  questions,
		Answers:    answers,
	})
}

type GetTestByIdResponse struct {
	Test      models.TestResponse          `json:"test"`
	Questions []models.GetQuestionResponse `json:"questions"`
	Answers   []models.QuestionAnswer      `json:"answers"`
}

type GetTestByIdResponse2 struct {
	Test       models.TestResponse           `json:"test"`
	TestAnswer models.TestAnswer             `json:"test_answer"`
	Questions  []models.GetQuestionResponse2 `json:"questions"`
	Answers    []models.QuestionAnswer       `json:"answers"`
}

// updateTest godoc
//
// @Summary     Update a test
// @Security    ApiKeyAuth
// @Description Update a test
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       id  path     int            true "A unique integer value identifying this test."
// @Success     200 {object} statusResponse "todo"
// @Failure     400 {object} errorResponse  "todo"
// @Failure     404 {object} errorResponse  "incorrect id"
// @Failure     500 {object} errorResponse  "todo"
// @Router      /test/{id} [put]
func (h *Handler) updateTest(c *gin.Context) {
	testId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, c.Param("id")+" not found", nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.UpdateTestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if input.Title == nil &&
		input.Description == nil &&
		input.RandomQuestionsOrder == nil &&
		input.QuestionsVisibility == nil &&
		input.Start == nil &&
		input.End == nil &&
		input.Duration_sec == nil {
		newErrorResponse(c, http.StatusBadRequest, "empty input", nil)
		return
	}

	err = h.services.Test.Update(userId, testId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// deleteTest godoc
//
// @Summary     Delete a test
// @Security    ApiKeyAuth
// @Description Delete a test
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       id  path     int            true "A unique integer value identifying this test."
// @Success     200 {object} statusResponse "todo"
// @Failure     404 {object} errorResponse  "incorrect id"
// @Failure     500 {object} errorResponse  "todo"
// @Router      /test/{id} [delete]
func (h *Handler) deleteTest(c *gin.Context) {
	testId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, c.Param("id")+" not found", nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	err = h.services.Test.Delete(userId, testId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// completeTest godoc
//
// @Summary     Complete a test
// @Security    ApiKeyAuth
// @Description Complete a test
// @Tags        test
// @Accept      json
// @Produce     json
// @Param       id  path     int            true "A unique integer value identifying this test."
// @Success     200 {object} statusResponse "todo"
// @Failure     404 {object} errorResponse  "incorrect id"
// @Failure     500 {object} errorResponse  "todo"
// @Router      /test-answer [post]
func (h *Handler) completeTest(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.TestAnswer
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	input.UserId = userId
	input.CompleteTime = time.Now()
	err = h.services.TestAnswer.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
