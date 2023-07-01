package repositories

import (
	"chat-app/models"
	"context"
	"gorm.io/gorm"
)

type FriendRepository interface {
	Create(ctx context.Context, friend models.Friend) error
	Delete(ctx context.Context, friend models.Friend) error
	GetAllByUserId(ctx context.Context, userId string) ([]models.Friend, error)
}

type friendRepository struct {
	database *gorm.DB
}

func (f *friendRepository) Create(ctx context.Context, friend models.Friend) error {
	var (
		db = f.database.WithContext(ctx)
	)
	return db.Model(&models.Friend{}).Create(&friend).Error
}

func (f *friendRepository) Delete(ctx context.Context, friend models.Friend) error {
	return f.database.WithContext(ctx).Model(&models.Friend{}).
		Where("user_id_1 = ? and user_id_2 = ?", friend.UserId1, friend.UserId2).
		Delete(&models.Friend{}).
		Error
}

func (f *friendRepository) GetAllByUserId(
	ctx context.Context,
	userId string,
) ([]models.Friend, error) {
	friends := make([]models.Friend, 0)
	err := f.database.WithContext(ctx).
		Model(&models.Friend{}).
		Where("user_id_1 = ? or user_id_2 = ?", userId, userId).
		Find(&friends).
		Error
	return friends, err
}

func NewFriendRepository(database *gorm.DB) FriendRepository {
	return &friendRepository{
		database: database,
	}
}
