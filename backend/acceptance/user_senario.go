package acceptance

import (
	"testing"
	"yangdongju/gtd_todo/internal/user"

	"github.com/stretchr/testify/assert"
)

func SignUpSenario(t *testing.T, userDsl userDsl) {
	payload := user.SignUpRequest{
		Email: "example@test.com",
		Password: "examplePasswords",
	}
	respond, err := userDsl.signUp(payload)

	assert.NoError(t, err)
	assert.Equal(t, respond.Email, payload.Email)
	assert.Equal(t, respond.ID, 1)
}
