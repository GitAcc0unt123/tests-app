package repository

import (
	"tests_app/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshSessionPostgres struct {
	db *sqlx.DB
}

func NewRefreshSessionPostgres(db *sqlx.DB) *RefreshSessionPostgres {
	return &RefreshSessionPostgres{db: db}
}

func (r *RefreshSessionPostgres) Create(refreshSession models.RefreshSession) error {
	row := r.db.QueryRow(`
	INSERT INTO refresh_sessions
	VALUES ($1, $2, $3, $4, $5, $6)`,
		refreshSession.RefreshToken,
		refreshSession.UserId,
		refreshSession.UserAgent,
		refreshSession.Fingerprint,
		refreshSession.Ip,
		refreshSession.ExpiresAt)

	return row.Err()
}

func (r *RefreshSessionPostgres) Get(refreshToken uuid.UUID) (models.RefreshSession, error) {
	var refreshSession models.RefreshSession
	err := r.db.Get(&refreshSession, `SELECT * FROM refresh_sessions WHERE refresh_token = $1`, refreshToken)
	return refreshSession, err
}

func (r *RefreshSessionPostgres) Revoke(refreshToken uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM refresh_sessions WHERE refresh_token = $1`, refreshToken)
	return err
}
