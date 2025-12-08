package user

import (
	"fmt"
	"net/http"
)

type userAlreadyExistsError struct {
	Code      int
	Message   string
	NestedErr error
}

func (e userAlreadyExistsError) Error() string {
	return e.Message
}

func newUserAlreadyExistsError(id int, email string) *userAlreadyExistsError {
	return &userAlreadyExistsError{
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

func newInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{
		Code:      http.StatusUnauthorized,
		Message:   "Invalid email or password",
		NestedErr: nil,
	}
}
