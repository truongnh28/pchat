package repositories

import (
	"chat-app/models"
	"context"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, req *models.File) error
	Get(ctx context.Context, id uint) (models.File, error)
	Delete(ctx context.Context, id uint) error
}

type fileRepository struct {
	database *gorm.DB
}

func (f *fileRepository) Delete(ctx context.Context, id uint) error {
	return f.database.WithContext(ctx).
		Model(models.File{}).
		Where("id = ?", id).
		Delete(models.File{}).
		Error
}

func (f *fileRepository) Get(ctx context.Context, id uint) (models.File, error) {
	file := models.File{}
	err := f.database.WithContext(ctx).
		Model(models.File{}).
		Where("id = ?", id).
		Find(&file).
		Error
	return file, err
}

func (f *fileRepository) Create(ctx context.Context, file *models.File) error {
	return f.database.WithContext(ctx).Model(models.File{}).Create(&file).Error
}

func NewFileRepository(database *gorm.DB) FileRepository {
	return &fileRepository{
		database: database,
	}
}
