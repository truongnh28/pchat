package repositories

import (
	"chat-app/internal/domain"
	"chat-app/models"
	"context"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, req models.Room) error
	Get(ctx context.Context, req domain.Room) ([]*models.Room, error)
	Update(ctx context.Context, req models.Room) error
	Delete(ctx context.Context, req models.Room) error
}

type roomRepository struct {
	database *gorm.DB
}

func (g *roomRepository) Delete(ctx context.Context, req models.Room) error {
	//TODO implement me
	panic("implement me")
}

func (g *roomRepository) Update(ctx context.Context, req models.Room) error {
	//TODO implement me
	panic("implement me")
}

func (g *roomRepository) Get(ctx context.Context, req domain.Room) ([]*models.Room, error) {
	room := make([]*models.Room, 0)
	db := g.database.WithContext(ctx).
		Model(models.Room{})
	if req.GroupId != "" {
		db = db.Where("group_id = ?", req.GroupId)
	}
	if req.UserId != "" {
		db = db.Where("user_id = ?", req.UserId)
	}
	err := db.Find(&room).Error
	return room, err
}

func (g *roomRepository) Create(ctx context.Context, room models.Room) error {
	var (
		db = g.database.WithContext(ctx)
	)
	return db.Model(models.Room{}).Create(&room).Error
}

func NewRoomRepository(database *gorm.DB) RoomRepository {
	return &roomRepository{
		database: database,
	}
}
