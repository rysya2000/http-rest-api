package store

import "errors"

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrNotAuthenticated = errors.New("not authenticated")

	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
)
