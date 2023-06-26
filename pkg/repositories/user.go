package repositories

import (
	"chat-app/models"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, req models.User) error
	FindByUserName(ctx context.Context, username string) (*models.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	UpdatePassword(
		ctx context.Context,
		req *chat_app.UpdateUserRequest,
	) (int64, error)
	UpdateStatus(
		ctx context.Context,
		email string,
		status models.AccountStatus,
	) (int64, error)
	Validate(ctx context.Context, req models.User) error
}

type userRepository struct {
	database *gorm.DB
}

func (a *userRepository) UpdateStatus(
	ctx context.Context,
	email string,
	status models.AccountStatus,
) (int64, error) {
	db := a.database.WithContext(ctx)
	result := db.Model(models.User{}).
		Select("status").
		Where("email = ?", email).
		Updates(models.User{Status: status})
	return result.RowsAffected, result.Error
}

func (a *userRepository) Validate(
	ctx context.Context,
	req models.User,
) error {
	var count = int64(0)
	err := a.database.WithContext(ctx).
		Model(models.User{}).
		Where("user_name = ? or email = ? or phone_number = ?", req.UserName, req.Email, req.PhoneNumber).
		Count(&count).
		Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username or email is exist")
	}
	return nil
}

func (a *userRepository) Get(ctx context.Context, userId string) (*models.User, error) {
	userProfiles := &models.User{}
	err := a.database.WithContext(ctx).
		Model(models.User{}).
		Where("user_id = ?", userId).
		First(&userProfiles).
		Error
	return userProfiles, err
}

func (a *userRepository) UpdatePassword(
	ctx context.Context,
	req *chat_app.UpdateUserRequest,
) (int64, error) {
	db := a.database.WithContext(ctx)
	result := db.Model(models.User{}).
		Select("password").
		Where("user_name = ?", req.GetUsername()).
		Updates(models.User{Password: req.GetPassword()})
	return result.RowsAffected, result.Error
}

func (a *userRepository) FindByUserName(
	ctx context.Context,
	username string,
) (*models.User, error) {
	userProfiles := &models.User{}
	err := a.database.WithContext(ctx).
		Model(models.User{}).
		Where("user_name = ?", username).
		Find(&userProfiles).
		Error
	return userProfiles, err
}

func (a *userRepository) FindByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (*models.User, error) {
	userProfiles := &models.User{}
	err := a.database.WithContext(ctx).
		Model(models.User{}).
		Where("phone_number = ?", phoneNumber).
		Find(&userProfiles).
		Error
	return userProfiles, err
}

func (a *userRepository) Create(ctx context.Context, user models.User) error {
	var (
		db = a.database.WithContext(ctx)
	)
	return db.Model(models.User{}).Create(&user).Error
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{
		database: database,
	}
}
