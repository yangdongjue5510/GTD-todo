package user

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============ Mock SignUpUsecase ============
type mockSignUpUsecase struct {
	signUpFuncStub func(request SignUpRequest) (*SignUpResponse, error)
}

func (m *mockSignUpUsecase) signUp(request SignUpRequest) (*SignUpResponse, error) {
	if m.signUpFuncStub != nil {
		return m.signUpFuncStub(request)
	}
	return nil, errors.New("signUpFuncStub not implemented")
}

// ============ Test Cases ============

func TestHandleSignUp_Success(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request SignUpRequest) (*SignUpResponse, error) {
			return &SignUpResponse{
				ID:    1,
				Email: request.Email,
			}, nil
		},
	}

	handler := newUserHandler(mockUsecase)
	request := SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	signUpRes, ok := res.(*SignUpResponse)

	// then
	assert.Equal(t, http.StatusCreated, code)
	assert.True(t, ok, "Expected *SignUpResponse type")
	assert.NotNil(t, signUpRes)
	assert.Equal(t, 1, signUpRes.ID)
	assert.Equal(t, "test@example.com", signUpRes.Email)
}

func TestHandleSignUp_UserAlreadyExists(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request SignUpRequest) (*SignUpResponse, error) {
			return nil, newUserAlreadyExistsError(1, request.Email)
		},
	}

	handler := newUserHandler(mockUsecase)
	request := SignUpRequest{
		Email:    "existing@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	errRes, ok := res.(ErrorResponse)

	// then
	assert.Equal(t, http.StatusBadRequest, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "User already exists")
	assert.Contains(t, errRes.Error, "existing@example.com")
}

func TestHandleSignUp_InternalServerError(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request SignUpRequest) (*SignUpResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := newUserHandler(mockUsecase)
	request := SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	errRes, ok := res.(ErrorResponse)

	// then
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "database connection failed")
}
