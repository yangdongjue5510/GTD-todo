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

type mockLoginUsecase struct {
	loginFuncStub func(request LoginRequest) (*LoginResponse, error)
}

func (m *mockLoginUsecase) Login(request LoginRequest) (*LoginResponse, error) {
	if m.loginFuncStub != nil {
		return m.loginFuncStub(request)
	}
	return nil, errors.New("loginStub not implemented")
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

	handler := &UserHandler{signUpUsecase: mockUsecase}
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

	handler := &UserHandler{signUpUsecase: mockUsecase}
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

	handler := &UserHandler{signUpUsecase: mockUsecase}
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

func TestHandleLogIn_Success(t *testing.T) {
	// given
	mockLoginUsecase := mockLoginUsecase{
		loginFuncStub: func(request LoginRequest) (*LoginResponse, error) {
			return &LoginResponse{"example_token"}, nil
		},
	}

	handler := &UserHandler{loginUsecase: &mockLoginUsecase}
	request := LoginRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	logInResponse, ok := res.(*LoginResponse)

	// then
	assert.Equal(t, http.StatusOK, code)
	assert.True(t, ok, "Expected type")
	assert.NotNil(t, logInResponse)
	assert.Equal(t, "example_token", logInResponse.Token)
}

func TestHandleLogIn_InvalidCredentials(t *testing.T) {
	// given
	mockLoginUsecase := &mockLoginUsecase{
		loginFuncStub: func(request LoginRequest) (*LoginResponse, error) {
			return nil, newInvalidCredentialsError()
		},
	}

	handler := &UserHandler{loginUsecase: mockLoginUsecase}
	request := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	errRes, ok := res.(ErrorResponse)

	// then
	assert.Equal(t, http.StatusUnauthorized, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Equal(t, "Invalid email or password", errRes.Error)
}

func TestHandleLogIn_InternalServerError(t *testing.T) {
	// given
	mockLoginUsecase := &mockLoginUsecase{
		loginFuncStub: func(request LoginRequest) (*LoginResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := &UserHandler{loginUsecase: mockLoginUsecase}
	request := LoginRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	errRes, ok := res.(ErrorResponse)

	// then
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "database connection failed")
}
