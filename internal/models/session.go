package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	//Id           int       `json:"id"            db:"id"`
	RefreshToken uuid.UUID `json:"refresh_token" db:"refresh_token"`
	UserId       int       `json:"user_id"       db:"user_id"`
	UserAgent    string    `json:"user_agent"    db:"user_agent"`
	Fingerprint  string    `json:"fingerprint"   db:"fingerprint"`
	Ip           net.IP    `json:"ip"            db:"ip"`
	ExpiresAt    time.Time `json:"expires_at"    db:"expires_at"`
}

func (r *RefreshSession) Expired() bool {
	return time.Now().After(r.ExpiresAt)
}
