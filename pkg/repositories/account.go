package repositories

import (
	"chat-app/models"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, req models.Account) error
	FindByUserName(ctx context.Context, username string) (*models.Account, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Account, error)
	Get(ctx context.Context, id uint64) (*models.Account, error)
	UpdatePassword(
		ctx context.Context,
		req *chat_app.UpdateAccountRequest,
	) (int64, error)
	UpdateStatus(
		ctx context.Context,
		email string,
		status models.AccountStatus,
	) (int64, error)
	Validate(ctx context.Context, req models.Account) error
}

type accountRepository struct {
	database *gorm.DB
}

func (a *accountRepository) UpdateStatus(
	ctx context.Context,
	email string,
	status models.AccountStatus,
) (int64, error) {
	db := a.database.WithContext(ctx)
	result := db.Model(models.Account{}).
		Select("status").
		Where("email = ?", email).
		Updates(models.Account{Status: status})
	return result.RowsAffected, result.Error
}

func (a *accountRepository) Validate(
	ctx context.Context,
	req models.Account,
) error {
	var count = int64(0)
	err := a.database.WithContext(ctx).
		Model(models.Account{}).
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

func (a *accountRepository) Get(ctx context.Context, id uint64) (*models.Account, error) {
	userProfiles := &models.Account{}
	err := a.database.WithContext(ctx).
		Model(models.Account{}).
		Where("id = ?", id).
		Find(&userProfiles).
		Error
	return userProfiles, err
}

func (a *accountRepository) UpdatePassword(
	ctx context.Context,
	req *chat_app.UpdateAccountRequest,
) (int64, error) {
	db := a.database.WithContext(ctx)
	result := db.Model(models.Account{}).
		Select("password").
		Where("user_name = ?", req.GetUsername()).
		Updates(models.Account{Password: req.GetPassword()})
	return result.RowsAffected, result.Error
}

func (a *accountRepository) FindByUserName(
	ctx context.Context,
	username string,
) (*models.Account, error) {
	userProfiles := &models.Account{}
	err := a.database.WithContext(ctx).
		Model(models.Account{}).
		Where("user_name = ?", username).
		Find(&userProfiles).
		Error
	return userProfiles, err
}

func (a *accountRepository) FindByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (*models.Account, error) {
	userProfiles := &models.Account{}
	err := a.database.WithContext(ctx).
		Model(models.Account{}).
		Where("phone_number = ?", phoneNumber).
		Find(&userProfiles).
		Error
	return userProfiles, err
}

func (a *accountRepository) Create(ctx context.Context, account models.Account) error {
	var (
		db = a.database.WithContext(ctx)
	)
	return db.Model(models.Account{}).Create(&account).Error
}

func NewAccountRepository(database *gorm.DB) AccountRepository {
	return &accountRepository{
		database: database,
	}
}
