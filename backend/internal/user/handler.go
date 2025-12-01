package user

import (
	"errors"
	"net/http"
)

type UserHandler struct {
	signUpUsecase signUpUsecase
	loginUsecase  loginUsecase
}

func (h *UserHandler) HandleSignUp(req SignUpRequest) (int, any) {
	res, err := h.signUpUsecase.signUp(req)
	if err != nil {
		return handleError(err)
	}
	return http.StatusCreated, res
}

func (h *UserHandler) HandleLogin(req loginRequest) (int, any) {
	res, err := h.loginUsecase.login(req)
	if err != nil {
		return handleError(err)
	}
	return http.StatusOK, res
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(err error) (int, ErrorResponse) {
	if err == nil {
		return 0, ErrorResponse{}
	}

	var userAlreadyExistsError *userAlreadyExistsError
	var invalidCredentialsError *invalidCredentialsError

	switch {
	case errors.As(err, &userAlreadyExistsError):
		return http.StatusBadRequest, ErrorResponse{Error: err.Error()}
	case errors.As(err, &invalidCredentialsError):
		return http.StatusUnauthorized, ErrorResponse{Error: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorResponse{Error: err.Error()}
	}
}
