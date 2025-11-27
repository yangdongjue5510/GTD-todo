package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type signUpUsecase interface {
	signUp(request SignUpRequest) (*SignUpResponse, error)
}

type userService struct {
	userRepository userRepository
}

func newUserService(repository userRepository) *userService {
	return &userService{userRepository: repository}
}

func (s *userService) signUp(req SignUpRequest) (*SignUpResponse, error) {
	foundUser, err := s.userRepository.findUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if foundUser != nil {
		return nil, newUserAlreadyExistsError(foundUser.ID, foundUser.Email)
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

	savedUser, err := s.userRepository.save(&newUser)
	return &SignUpResponse{
		ID:    savedUser.ID,
		Email: savedUser.Email,
	}, err
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
