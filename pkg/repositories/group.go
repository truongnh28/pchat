package repositories

import (
	"chat-app/models"
	"context"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Create(ctx context.Context, req models.Group) error
	Get(ctx context.Context, id string) (*models.Group, error)
	Update(ctx context.Context, req models.Group) error
	Delete(ctx context.Context, id string) error
}

type groupRepository struct {
	database *gorm.DB
}

func (g *groupRepository) Delete(ctx context.Context, groupId string) (err error) {
	tx := g.database.WithContext(ctx).Begin()
	defer func() {
		err = finalizeTransaction(tx, err)
	}()
	err = tx.Model(&models.Room{}).
		Where("group_id = ?", groupId).
		Delete(&models.Room{}).
		Error
	if err != nil {
		return
	}
	err = tx.Model(&models.Group{}).
		Where("group_id = ?", groupId).
		Delete(&models.Group{}).
		Error

	return
}

func (g *groupRepository) Update(ctx context.Context, req models.Group) (err error) {
	tx := g.database.WithContext(ctx).Begin()
	defer func() {
		err = finalizeTransaction(tx, err)
	}()
	params := make(map[string]any)
	if req.Name != "" {
		params["name"] = req.Name
	}
	if req.FileId != 0 {
		params["file_id"] = req.FileId
	}
	err = tx.Model(&models.Group{}).
		Where("group_id = ?", req.GroupId).
		Updates(params).Error
	if err != nil {
		return
	}

	return
}

func (g *groupRepository) Get(ctx context.Context, id string) (*models.Group, error) {
	group := &models.Group{}
	err := g.database.WithContext(ctx).
		Model(&models.Group{}).
		Where("group_id = ?", id).
		Find(&group).
		Error
	return group, err
}

func (g *groupRepository) Create(ctx context.Context, group models.Group) error {
	var (
		db = g.database.WithContext(ctx)
	)
	return db.Model(&models.Group{}).Create(&group).Error
}

func NewGroupRepository(database *gorm.DB) GroupRepository {
	return &groupRepository{
		database: database,
	}
}
