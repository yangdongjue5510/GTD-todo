package user

import (
	"errors"
	"net/http"
)

type UserHandler struct {
	signUpUsecase signUpUsecase
}

func newUserHandler(signUpUsecase signUpUsecase) *UserHandler {
	return &UserHandler{
		signUpUsecase: signUpUsecase,
	}
}

func (h *UserHandler) HandleSignUp(req SignUpRequest) (int, any) {
	res, err := h.signUpUsecase.signUp(req)
	if err != nil {
		return handleError(err)
	}
	return http.StatusCreated, res
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(err error) (int, ErrorResponse) {
	if err == nil {
		return 0, ErrorResponse{}
	}

	var userAlreadyExistsError *userAlreadyExistsError

	switch {
	case errors.As(err, &userAlreadyExistsError):
		return http.StatusBadRequest, ErrorResponse{Error: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorResponse{Error: err.Error()}
	}
}
