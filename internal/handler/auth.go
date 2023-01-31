package handler

import (
	"net"
	"net/http"
	"tests_app/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// signUp godoc
//
// @Summary     Create an account
// @Description Create an account
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       data body     models.User true "account info"
// @Success     201  {object} idResponse  "Account id"
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body", err)
		return
	}

	id, err := h.services.User.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	c.JSON(http.StatusCreated, idResponse{
		Id: id,
	})
}

type signInInput struct {
	Username    string `json:"username"    binding:"required,min=3,max=255"`
	Password    string `json:"password"    binding:"required,min=8,max=255"`
	Fingerprint string `json:"fingerprint" binding:"required,max=200"`
}

// signIn godoc
//
// @Summary     Log in account
// @Description Log in account
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       data body     signInInput true "sign in info"
// @Success     200  {object} tokenResponse
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body", err)
		return
	}

	user, err := h.services.User.Get(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	accessToken, err := h.services.TokenManager.NewJWT(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	newRefreshToken, err := h.services.TokenManager.NewRefreshToken()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	refreshSession := models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserId:       user.Id,
		UserAgent:    c.Request.UserAgent(),
		Fingerprint:  input.Fingerprint,
		Ip:           net.IP(c.ClientIP()),
		ExpiresAt:    time.Now().Add(h.config.JWT.RefreshTokenTTL),
	}

	err = h.services.RefreshSession.Create(refreshSession)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	// browser
	c.SetCookie(
		"refresh_token",
		refreshSession.RefreshToken.String(),
		int(h.config.JWT.RefreshTokenTTL.Seconds()),
		"/api/auth",
		"",
		false,
		true)
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshSession.RefreshToken.String(), // mobile app
	})
}

func (h *Handler) signOut(c *gin.Context) {
	refreshSessionCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "", err)
		return
	}

	refreshToken, err := uuid.Parse(refreshSessionCookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "", err)
		return
	}

	if err := h.services.RefreshSession.Revoke(refreshToken); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", err)
		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/api/auth",
		"",
		false,
		true)
	c.Status(http.StatusOK)
}

// RefreshSessions godoc
//
// @Summary     Refresh access token
// @Description Refresh access token
// @Tags        account
// @Accept      json
// @Produce     json
// @Success     200 {object} tokenResponse
// @Failure     400 {object} errorResponse
// @Failure     500 {object} errorResponse
// @Router      /auth/refresh-token [post]
func (h *Handler) refreshToken(c *gin.Context) {
	refreshTokenCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no cookie", err)
		return
	}

	var input struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body", err)
		return
	}

	refreshToken, err := uuid.Parse(refreshTokenCookie.Value)
	if err != nil {
		c.SetCookie("refresh_token", "", 0, "", "", false, false)
		newErrorResponse(c, http.StatusBadRequest, "invalid uuid cookie", err)
		return
	}

	refreshSession, err := h.services.RefreshSession.Get(refreshToken)
	if err != nil {
		c.SetCookie("refresh_token", "", -1, "", "", false, false)
		newErrorResponse(c, http.StatusInternalServerError, "uuid not found", err)
		return
	}

	if refreshSession.Expired() {
		c.SetCookie("refresh_token", "", -1, "", "", false, false)
		newErrorResponse(c, http.StatusInternalServerError, "uuid expired", err)
		return
	}

	// проверить клиент
	if refreshSession.UserAgent != c.Request.UserAgent() ||
		!refreshSession.Ip.Equal(net.IP(c.ClientIP())) ||
		refreshSession.Fingerprint != input.Fingerprint {
		c.SetCookie("refresh_token", "", -1, "", "", false, false)
		newErrorResponse(c, http.StatusInternalServerError, "uuid expired", err)
		return
	}

	accessToken, err := h.services.TokenManager.NewJWT(refreshSession.UserId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	newRefreshToken, err := h.services.TokenManager.NewRefreshToken()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	newRefreshSession := models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserId:       refreshSession.UserId,
		UserAgent:    refreshSession.UserAgent,
		Fingerprint:  refreshSession.Fingerprint,
		Ip:           refreshSession.Ip,
		ExpiresAt:    time.Now().Add(h.config.JWT.RefreshTokenTTL),
	}

	err = h.services.RefreshSession.Create(newRefreshSession)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	if err := h.services.RefreshSession.Revoke(refreshToken); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", err)
		return
	}

	// browser
	c.SetCookie(
		"refresh_token",
		newRefreshSession.RefreshToken.String(),
		int(h.config.JWT.RefreshTokenTTL.Seconds()),
		"/api/auth",
		"",
		false,
		true)
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshSession.RefreshToken.String(), // mobile app
	})
}

// updateUser godoc
//
// @Summary     Update account info
// @Security    ApiKeyAuth
// @Description Update account info
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       data body     models.UpdateUserInput true "update info"
// @Success     200  {object} statusResponse
// @Failure     400  {object} errorResponse
// @Failure     500  {object} errorResponse
// @Router      /auth/update-user [put]
func (h *Handler) updateUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if input.Name == nil &&
		input.Email == nil &&
		input.Username == nil &&
		input.Password == nil {
		newErrorResponse(c, http.StatusBadRequest, "empty input", nil)
		return
	}

	err = h.services.User.Update(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
