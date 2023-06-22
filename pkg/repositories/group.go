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
}

type groupRepository struct {
	database *gorm.DB
}

func (g *groupRepository) Update(ctx context.Context, req models.Group) error {
	//TODO implement me
	panic("implement me")
}

func (g *groupRepository) Get(ctx context.Context, id string) (*models.Group, error) {
	group := &models.Group{}
	err := g.database.WithContext(ctx).
		Model(models.Group{}).
		Where("id = ?", id).
		Find(&group).
		Error
	return group, err
}

func (g *groupRepository) Create(ctx context.Context, group models.Group) error {
	var (
		db = g.database.WithContext(ctx)
	)
	return db.Model(models.Group{}).Create(&group).Error
}

func NewGroupRepository(database *gorm.DB) GroupRepository {
	return &groupRepository{
		database: database,
	}
}
