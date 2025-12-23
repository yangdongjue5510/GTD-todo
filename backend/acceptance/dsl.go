package acceptance

import (
	"encoding/json"
	"errors"
	"net/http"

	"yangdongju/gtd_todo/internal/user"
)

type userDsl interface {
	signUp(user.SignUpRequest) (user.SignUpResponse, error)
}

type userDslImpl struct {
	apiDriver apiDriver
}

func (u userDslImpl) signUp(payload user.SignUpRequest) (user.SignUpResponse, error) {
	apiResp, err := u.apiDriver.call(http.MethodPost, "/api/auth/signup", payload, nil)
	if err != nil {
		return user.SignUpResponse{}, err
	}

	if len(apiResp.Body) > 0 {
		var parsed user.SignUpResponse
		if err := json.Unmarshal(apiResp.Body, &parsed); err != nil {
			return user.SignUpResponse{}, err
		}
		return parsed, nil
	}

	return user.SignUpResponse{}, errors.New("response body is empty.")
}
