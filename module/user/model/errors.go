package model

import "errors"

var (
	ErrUserDataIsRequired       = errors.New("user data is required")
	ErrUserUpdateDataIsRequired = errors.New("user update data is required")
	ErrUserIDIsBlank            = errors.New("user id is required")
	ErrUserEmailIsBlank         = errors.New("user email is required")
	ErrUserPasswordHashIsBlank  = errors.New("user password hash is required")
	ErrUserFullNameIsBlank      = errors.New("user full name is required")
	ErrUserRoleIsBlank          = errors.New("user role is required")
	ErrUserEmailAlreadyExists   = errors.New("user email already exists")
	ErrUserNotFound             = errors.New("user not found")
)
