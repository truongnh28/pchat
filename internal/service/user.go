package service

import (
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type UserService interface {
	GetByUserId(userID string) (domain.UserDetail, error)
	UpdateUserOnlineStatusByUserID(userId string, isOnline bool) error
	GetAllOnlineUsers(userId string) ([]domain.UserDetail, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func (u *userServiceImpl) GetAllOnlineUsers(userId string) ([]domain.UserDetail, error) {
	return nil, nil
}

func (u *userServiceImpl) GetByUserId(userID string) (domain.UserDetail, error) {
	return domain.UserDetail{
		ID:          "",
		Username:    "",
		PhoneNumber: "",
		Password:    "",
		Online:      false,
		SocketId:    "",
	}, nil
}

func (u *userServiceImpl) UpdateUserOnlineStatusByUserID(userId string, isOnline bool) error {
	return nil
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}
