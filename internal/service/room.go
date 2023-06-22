package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/models"
	"chat-app/pkg/repositories"
	"context"
	"github.com/golang/glog"
	"github.com/whatvn/denny"
)

type RoomService interface {
	Create(
		ctx context.Context,
		req domain.Room,
	) common.SubReturnCode
	Get(
		ctx context.Context,
		req domain.Room,
	) ([]domain.Room, common.SubReturnCode)
}

type roomService struct {
	roomRepository repositories.RoomRepository
}

func (a *roomService) Get(
	ctx context.Context,
	req domain.Room,
) ([]domain.Room, common.SubReturnCode) {
	var (
		rooms  = make([]*models.Room, 0)
		err    = error(nil)
		resp   = make([]domain.Room, 0)
		logger = denny.GetLogger(ctx)
	)
	if req.GroupId == "" || req.UserId == "" {
		return resp, common.SystemError
	}
	rooms, err = a.roomRepository.Get(ctx, req)
	if err != nil {
		glog.Errorf("Find Room fail: %s", err)
		logger.WithError(err)
		return resp, common.SystemError
	}
	for _, room := range rooms {
		resp = append(resp, domain.Room{
			GroupId: room.GroupId,
			UserId:  room.UserId,
		})
	}
	return resp, common.OK
}

func (a *roomService) Create(
	ctx context.Context,
	req domain.Room,
) common.SubReturnCode {
	var (
		logger = denny.GetLogger(ctx)
	)
	acc := models.Room{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	}

	err := a.roomRepository.Create(ctx, acc)
	if err != nil {
		logger.Errorln("Create Room service err: ", err)
		return common.SystemError
	}

	return common.OK
}

func NewRoomService(
	roomRepository repositories.RoomRepository,
) RoomService {
	return &roomService{
		roomRepository: roomRepository,
	}
}
