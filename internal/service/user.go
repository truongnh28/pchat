package service

import (
	"chat-app/pkg/repositories"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type UserService interface {
	//TODO: handle func
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}
