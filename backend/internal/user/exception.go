package user

import (
	"fmt"
	"net/http"
)

type UserAlreadyExistsError struct {
	Code      int
	Message   string
	NestedErr error
}

func (e UserAlreadyExistsError) Error() string {
	return e.Message
}

func NewUserAlreadyExistsError(id int, email string) *UserAlreadyExistsError {
	return &UserAlreadyExistsError{
		Code:      http.StatusBadRequest,
		Message:   fmt.Sprintf("User already exists. id=%v & email=%v", id, email),
		NestedErr: nil,
	}
}

type InvalidCredentialsError struct {
	Code      int
	Message   string
	NestedErr error
}

func (e InvalidCredentialsError) Error() string {
	return e.Message
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{
		Code:      http.StatusUnauthorized,
		Message:   "Invalid email or password",
		NestedErr: nil,
	}
}
