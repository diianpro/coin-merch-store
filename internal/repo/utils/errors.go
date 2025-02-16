package utils

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	ErrNotEnoughBalance = errors.New("not enough balance")
)
