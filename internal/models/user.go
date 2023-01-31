package models

import "time"

type User struct {
	Id          int        `json:"-"        db:"id"`
	Name        string     `json:"name"     db:"name"         binding:"max=255"`
	Email       string     `json:"email"    db:"email"        binding:"email,max=255"`
	Username    string     `json:"username" db:"username"     binding:"required,min=3,max=255"`
	Password    string     `json:"password" db:"password"     binding:"required,min=8,max=255"`
	CreatedAt   time.Time  `json:"-"        db:"created_at"`
	ActivatedAt *time.Time `json:"-"        db:"activated_at"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}
