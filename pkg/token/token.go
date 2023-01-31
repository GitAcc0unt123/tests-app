package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenManager interface {
	NewJWT(userId int) (string, error)
	ParseJWT(accessToken string) (int, error)
	NewRefreshToken() string
}

type Manager struct {
	accessTokenTTL time.Duration
	signingKey     string
}

func NewManager(accessTokenTTL time.Duration, signingKey string) *Manager {
	return &Manager{
		accessTokenTTL: accessTokenTTL,
		signingKey:     signingKey,
	}
}

func (m *Manager) NewJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
		Subject:   strconv.Itoa(userId),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseJWT(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token == nil {
			return nil, errors.New("token is nil")
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *jwt.RegisteredClaims")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("sub field is not int")
	}

	return userId, nil
}

func (m *Manager) NewRefreshToken() (uuid.UUID, error) {
	return uuid.NewRandom()
}
