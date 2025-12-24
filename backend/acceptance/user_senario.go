package acceptance

import (
	"testing"
	"yangdongju/gtd_todo/internal/user"
	"yangdongju/gtd_todo/testhelper"

	"github.com/stretchr/testify/assert"
)

func signUpSenario(t *testing.T, userDsl userDsl) {
	testhelper.CleanUp()
	payload := user.SignUpRequest{
		Email:    "example@test.com",
		Password: "examplePasswords",
	}
	respond, err := userDsl.signUp(payload)

	assert.NoError(t, err)
	assert.Equal(t, respond.Email, payload.Email)
	assert.Equal(t, respond.ID, 1)
}

func loginSenario(t *testing.T, userDsl userDsl) {
	testhelper.CleanUp()
	signUpRequest := user.SignUpRequest{
		Email:    "example@test.com",
		Password: "examplePasswords",
	}
	_,_ = userDsl.signUp(signUpRequest)

	loginRequest := user.LoginRequest(signUpRequest)
	loginResponse, err := userDsl.login(loginRequest)

	assert.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Token)
}
