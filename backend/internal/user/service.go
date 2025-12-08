package user

type userService struct {
	userRepository UserRepository
	tokenIssuer    Issuer
	tokenParser    Parser
}

func NewUserService(repository UserRepository, tokenIssuer Issuer, tokenParser Parser) *userService {
	return &userService{
		userRepository: repository,
		tokenIssuer:    tokenIssuer,
		tokenParser:    tokenParser,
	}
}
