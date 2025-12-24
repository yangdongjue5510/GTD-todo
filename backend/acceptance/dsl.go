package acceptance

import (
	"encoding/json"
	"errors"
	"net/http"

	"yangdongju/gtd_todo/internal/user"
)

type userDsl interface {
	signUp(user.SignUpRequest) (user.SignUpResponse, error)
	login(user.LoginRequest) (user.LoginResponse, error)
}

type userDslImpl struct {
	apiDriver apiDriver
}

func (u userDslImpl) signUp(payload user.SignUpRequest) (user.SignUpResponse, error) {
	apiResp, err := u.apiDriver.call(http.MethodPost, "/api/auth/signup", payload, nil)
	if err != nil {
		return user.SignUpResponse{}, err
	}

	return unmarshal(apiResp, &user.SignUpResponse{})
}

func (u userDslImpl) login(payload user.LoginRequest) (user.LoginResponse, error) {
	apiResp, err := u.apiDriver.call(http.MethodPost, "/api/auth/login", payload, nil)

	if err != nil {
		return user.LoginResponse{}, err
	}

	return unmarshal(apiResp, &user.LoginResponse{})
}

func unmarshal[T any](apiResponse apiResponse, target *T) (T, error) {
	if target == nil {
		return *new(T), errors.New("response target to unmarshal is nil.")
	}

	if len(apiResponse.Body) <= 0 {
		return *new(T),errors.New("response body is empty.")
	}

	if err := json.Unmarshal(apiResponse.Body, target); err != nil {
		return *new(T), err
	}
	return *target, nil
}
