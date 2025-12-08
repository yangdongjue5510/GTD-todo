package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SignUpUsecase interface {
	SignUp(request SignUpRequest) (*SignUpResponse, error)
}

func (s *userService) SignUp(req SignUpRequest) (*SignUpResponse, error) {
	foundUser, err := s.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if foundUser != nil {
		return nil, NewUserAlreadyExistsError(foundUser.ID, foundUser.Email)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	savedUser, err := s.userRepository.Save(&newUser)
	if err != nil {
		return nil, err
	}
	return &SignUpResponse{
		ID:    savedUser.ID,
		Email: savedUser.Email,
	}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SignUpResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
