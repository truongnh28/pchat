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
		request domain.Room,
	) common.SubReturnCode
	Get(
		ctx context.Context,
		req domain.Room,
	) ([]domain.Room, common.SubReturnCode)
}

type roomService struct {
	roomRepository repositories.RoomRepository
}

func (r *roomService) Get(
	ctx context.Context,
	req domain.Room,
) ([]domain.Room, common.SubReturnCode) {
	var (
		rooms  = make([]*models.Room, 0)
		err    = error(nil)
		resp   = make([]domain.Room, 0)
		logger = denny.GetLogger(ctx)
	)
	if req.GroupId == "" && req.UserId == "" {
		return resp, common.SystemError
	}
	rooms, err = r.roomRepository.Get(ctx, req)
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

func (r *roomService) Create(
	ctx context.Context,
	request domain.Room,
) common.SubReturnCode {
	var (
		err    = error(nil)
		logger = denny.GetLogger(ctx)
	)
	if request.GroupId == "" || request.UserId == "" {
		logger.WithError(err).Errorf("invalid request: %s", err)
		return common.InvalidRequest
	}
	err = r.roomRepository.Create(ctx, models.Room{
		GroupId: request.GroupId,
		UserId:  request.UserId,
	})
	if err != nil {
		logger.WithError(err).Errorf("create room fail: %s", err)
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
