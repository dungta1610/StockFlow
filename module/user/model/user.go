package model

import (
	"strings"
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserCreate struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
}

func (u *UserCreate) Validate() error {
	if u == nil {
		return ErrUserDataIsRequired
	}

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.PasswordHash = strings.TrimSpace(u.PasswordHash)
	u.FullName = strings.TrimSpace(u.FullName)
	u.Role = strings.TrimSpace(u.Role)

	if u.Email == "" {
		return ErrUserEmailIsBlank
	}

	if u.PasswordHash == "" {
		return ErrUserPasswordHashIsBlank
	}

	if u.FullName == "" {
		return ErrUserFullNameIsBlank
	}

	if u.Role == "" {
		return ErrUserRoleIsBlank
	}

	return nil
}

type UserUpdate struct {
	PasswordHash *string `json:"password_hash"`
	FullName     string  `json:"full_name"`
	Role         string  `json:"role"`
	IsActive     *bool   `json:"is_active"`
}

func (u *UserUpdate) Validate() error {
	if u == nil {
		return ErrUserUpdateDataIsRequired
	}

	u.FullName = strings.TrimSpace(u.FullName)
	u.Role = strings.TrimSpace(u.Role)

	if u.FullName == "" {
		return ErrUserFullNameIsBlank
	}

	if u.Role == "" {
		return ErrUserRoleIsBlank
	}

	if u.PasswordHash != nil {
		trimmed := strings.TrimSpace(*u.PasswordHash)
		if trimmed == "" {
			return ErrUserPasswordHashIsBlank
		}
		*u.PasswordHash = trimmed
	}

	return nil
}

type Filter struct {
	Email    string `json:"email" form:"email"`
	FullName string `json:"full_name" form:"full_name"`
	Role     string `json:"role" form:"role"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.Email = strings.TrimSpace(strings.ToLower(f.Email))
	f.FullName = strings.TrimSpace(f.FullName)
	f.Role = strings.TrimSpace(f.Role)
}
