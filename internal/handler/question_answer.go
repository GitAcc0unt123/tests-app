package handler

import (
	"net/http"
	"tests_app/internal/models"

	"github.com/gin-gonic/gin"
)

// getAnswer godoc
//
// @Summary     Get an answer to the question
// @Security    ApiKeyAuth
// @Description Get an answer to the question
// @Tags        answer
// @Accept      json
// @Produce     json
// @Param       input body     getAnswerInput true "question id"
// @Success     200   {object} models.QuestionAnswer
// @Failure     400   {object} errorResponse
// @Failure     403   {object} errorResponse
// @Failure     500   {object} errorResponse
// @Router      /answer [get]
func (h *Handler) getAnswer(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error(), err)
		return
	}

	var input getAnswerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	question, err := h.services.QuestionAnswer.GetAnswerByQuestionId(userId, input.QuestionId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, question)
}

type getAnswerInput struct {
	QuestionId int `json:"question_id" binding:"required"`
}

// updateAnswer godoc
//
// @Summary     Create or update an answer
// @Security    ApiKeyAuth
// @Description Create or update an answer
// @Tags        answer
// @Accept      json
// @Produce     json
// @Param       data body     models.UpsertQuestionAnswerInput true "answer info"
// @Success     200  {object} statusResponse
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /answer [post]
func (h *Handler) setAnswer(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.UpsertQuestionAnswerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	err = h.services.QuestionAnswer.Upsert(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
