package user

import (
	"errors"
	"net/http"
)

type UserHandler struct {
	signUpUsecase SignUpUsecase
	loginUsecase  LoginUsecase
}

func NewUserHandler(loginUsecase LoginUsecase, signUpUsecase SignUpUsecase) *UserHandler {
	return &UserHandler{
		loginUsecase:  loginUsecase,
		signUpUsecase: signUpUsecase,
	}
}

func (h *UserHandler) HandleSignUp(req SignUpRequest) (int, any) {
	res, err := h.signUpUsecase.SignUp(req)
	if err != nil {
		return handleError(err)
	}
	return http.StatusCreated, res
}

func (h *UserHandler) HandleLogin(req LoginRequest) (int, any) {
	res, err := h.loginUsecase.Login(req)
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

	var userAlreadyExistsError *UserAlreadyExistsError
	var invalidCredentialsError *InvalidCredentialsError

	switch {
	case errors.As(err, &userAlreadyExistsError):
		return http.StatusBadRequest, ErrorResponse{Error: err.Error()}
	case errors.As(err, &invalidCredentialsError):
		return http.StatusUnauthorized, ErrorResponse{Error: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorResponse{Error: err.Error()}
	}
}
