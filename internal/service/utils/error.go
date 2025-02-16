package utils

import (
	"github.com/pkg/errors"
)

var (
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotParseToken = errors.New("cannot parse token")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrUserNotFound      = errors.New("user not found")
	ErrCannotGetUser     = errors.New("cannot get user")

	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrCannotCreateAccount  = errors.New("cannot create account")
	ErrAccountNotFound      = errors.New("account not found")
	ErrCannotGetAccount     = errors.New("cannot get account")
)
