package repositories

import (
	"chat-app/internal/domain"
	"chat-app/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mocks

type UserRepository interface {
	Create(user *models.Account) error
	GetByUserId(userId string) (*domain.UserDetail, error)
	UpdateUserOnlineStatusByUserID(userID string, status bool) error
	GetUserByUsername(username string) (*domain.UserDetail, error)
	IsUsernameAvailableQueryHandler(username string) bool
	GetAllOnlineUsers(userID string) ([]*domain.UserDetailsRequest, error)
}

type userRepositoryImpl struct {
	database *gorm.DB
}

func (n *userRepositoryImpl) GetByUserId(userId string) (*domain.UserDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (n *userRepositoryImpl) UpdateUserOnlineStatusByUserID(userID string, status bool) error {
	//TODO implement me
	panic("implement me")
}

func (n *userRepositoryImpl) GetUserByUsername(username string) (*domain.UserDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (n *userRepositoryImpl) IsUsernameAvailableQueryHandler(username string) bool {
	//TODO implement me
	panic("implement me")
}

func (n *userRepositoryImpl) GetAllOnlineUsers(
	userID string,
) ([]*domain.UserDetailsRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (n *userRepositoryImpl) Create(
	user *models.Account,
) error {
	return n.database.
		Model(models.Account{}).
		Create(&user).
		Error
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		database: database,
	}
}
