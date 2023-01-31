package handler

import (
	"net/http"
	"strconv"
	"tests_app/internal/models"

	"github.com/gin-gonic/gin"
)

// createQuestion godoc
//
// @Summary     Create question
// @Security    ApiKeyAuth
// @Description Create question
// @Tags        question
// @Accept      json
// @Produce     json
// @Param       data body     models.Question true "question info"
// @Success     200  {object} idResponse
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /question [post]
func (h *Handler) createQuestion(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.Question
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	QuestionId, err := h.services.Question.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, idResponse{
		Id: QuestionId,
	})
}

// getQuestionById godoc
//
// @Summary     Get question
// @Security    ApiKeyAuth
// @Description Get question
// @Tags        question
// @Accept      json
// @Produce     json
// @Param       id  path     int true "A unique integer value identifying this question."
// @Success     200 {object} models.Question
// @Failure     404 {object} errorResponse
// @Failure     500 {object} errorResponse
// @Router      /question/{id} [get]
func (h *Handler) getQuestionById(c *gin.Context) {
	QuestionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, c.Param("id")+" not found", nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	Question, err := h.services.Question.GetById(userId, QuestionId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, Question)
}

// updateQuestion godoc
//
// @Summary     Update question
// @Security    ApiKeyAuth
// @Description Update question
// @Tags        question
// @Accept      json
// @Produce     json
// @Param       id  path     int            true "A unique integer value identifying this question."
// @Success     200 {object} statusResponse "todo"
// @Failure     400 {object} errorResponse  "todo"
// @Failure     404 {object} errorResponse  "incorrect id"
// @Failure     500 {object} errorResponse  "todo"
// @Router      /question/{id} [put]
func (h *Handler) updateQuestion(c *gin.Context) {
	QuestionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, c.Param("id")+" not found", nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.UpdateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if input.Text == nil &&
		input.AnswerType == nil &&
		input.ShowAnswers == nil &&
		input.TrueAnswers == nil {
		newErrorResponse(c, http.StatusBadRequest, "empty input", nil)
		return
	}

	err = h.services.Question.Update(userId, QuestionId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// deleteQuestion godoc
//
// @Summary     Delete question
// @Security    ApiKeyAuth
// @Description Delete question
// @Tags        question
// @Accept      json
// @Produce     json
// @Param       id  path     int true "A unique integer value identifying this question."
// @Success     200 {object} statusResponse
// @Failure     404 {object} errorResponse
// @Failure     500 {object} errorResponse
// @Router      /question/{id} [delete]
func (h *Handler) deleteQuestion(c *gin.Context) {
	QuestionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, c.Param("id")+" not found", nil)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	err = h.services.Question.Delete(userId, QuestionId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// deleteQuestion godoc
//
// @Summary     Delete question
// @Security    ApiKeyAuth
// @Description Delete question
// @Tags        question
// @Accept      json
// @Produce     json
// @Param       id  body     getNextQuestionInput true "Get next question"
// @Success     200 {object} models.GetQuestionResponse
// @Failure     404 {object} errorResponse
// @Failure     500 {object} errorResponse
// @Router      /question/next [get]
func (h *Handler) getNextQuestion(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input getNextQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	question, err := h.services.Question.GetNext(userId, input.TestId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, question)
}

type getNextQuestionInput struct {
	TestId int `json:"test_id" binding:"required"`
}
