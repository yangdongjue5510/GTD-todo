package user_test

import (
	"errors"
	"net/http"
	"testing"
	"yangdongju/gtd_todo/internal/user"

	"github.com/stretchr/testify/assert"
)

// ============ Mock SignUpUsecase ============
type mockSignUpUsecase struct {
	signUpFuncStub func(request user.SignUpRequest) (*user.SignUpResponse, error)
}

func (m *mockSignUpUsecase) SignUp(request user.SignUpRequest) (*user.SignUpResponse, error) {
	if m.signUpFuncStub != nil {
		return m.signUpFuncStub(request)
	}
	return nil, errors.New("signUpFuncStub not implemented")
}

type mockLoginUsecase struct {
	loginFuncStub func(request user.LoginRequest) (*user.LoginResponse, error)
}

func (m *mockLoginUsecase) Login(request user.LoginRequest) (*user.LoginResponse, error) {
	if m.loginFuncStub != nil {
		return m.loginFuncStub(request)
	}
	return nil, errors.New("loginStub not implemented")
}

// ============ Test Cases ============

func TestHandleSignUp_Success(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request user.SignUpRequest) (*user.SignUpResponse, error) {
			return &user.SignUpResponse{
				ID:    1,
				Email: request.Email,
			}, nil
		},
	}

	handler := user.NewUserHandler(nil, mockUsecase)
	request := user.SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	signUpRes, ok := res.(*user.SignUpResponse)

	// then
	assert.Equal(t, http.StatusCreated, code)
	assert.True(t, ok, "Expected *user.SignUpResponse type")
	assert.NotNil(t, signUpRes)
	assert.Equal(t, 1, signUpRes.ID)
	assert.Equal(t, "test@example.com", signUpRes.Email)
}

func TestHandleSignUp_UserAlreadyExists(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request user.SignUpRequest) (*user.SignUpResponse, error) {
			return nil, user.NewUserAlreadyExistsError(1, request.Email)
		},
	}

	handler := user.NewUserHandler(nil, mockUsecase)
	request := user.SignUpRequest{
		Email:    "existing@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	errRes, ok := res.(user.ErrorResponse)

	// then
	assert.Equal(t, http.StatusBadRequest, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "User already exists")
	assert.Contains(t, errRes.Error, "existing@example.com")
}

func TestHandleSignUp_InternalServerError(t *testing.T) {
	// given
	mockUsecase := &mockSignUpUsecase{
		signUpFuncStub: func(request user.SignUpRequest) (*user.SignUpResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}
	handler := user.NewUserHandler(nil, mockUsecase)
	request := user.SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	code, res := handler.HandleSignUp(request)
	errRes, ok := res.(user.ErrorResponse)

	// then
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "database connection failed")
}

func TestHandleLogIn_Success(t *testing.T) {
	// given
	mockLoginUsecase := mockLoginUsecase{
		loginFuncStub: func(request user.LoginRequest) (*user.LoginResponse, error) {
			return &user.LoginResponse{"example_token"}, nil
		},
	}

	handler := user.NewUserHandler(&mockLoginUsecase, nil)
	request := user.LoginRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	logInResponse, ok := res.(*user.LoginResponse)

	// then
	assert.Equal(t, http.StatusOK, code)
	assert.True(t, ok, "Expected type")
	assert.NotNil(t, logInResponse)
	assert.Equal(t, "example_token", logInResponse.Token)
}

func TestHandleLogIn_InvalidCredentials(t *testing.T) {
	// given
	mockLoginUsecase := &mockLoginUsecase{
		loginFuncStub: func(request user.LoginRequest) (*user.LoginResponse, error) {
			return nil, user.NewInvalidCredentialsError()
		},
	}

	handler := user.NewUserHandler(mockLoginUsecase, nil)
	request := user.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	errRes, ok := res.(user.ErrorResponse)

	// then
	assert.Equal(t, http.StatusUnauthorized, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Equal(t, "Invalid email or password", errRes.Error)
}

func TestHandleLogIn_InternalServerError(t *testing.T) {
	// given
	mockLoginUsecase := &mockLoginUsecase{
		loginFuncStub: func(request user.LoginRequest) (*user.LoginResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := user.NewUserHandler(mockLoginUsecase, nil)
	request := user.LoginRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// when
	code, res := handler.HandleLogin(request)
	errRes, ok := res.(user.ErrorResponse)

	// then
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.True(t, ok, "Expected ErrorResponse type")
	assert.Contains(t, errRes.Error, "database connection failed")
}
