package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	repositories2 "chat-app/internal/repositories"
	"chat-app/models"
	"context"
	"github.com/google/uuid"
	"github.com/whatvn/denny"
)

type GroupService interface {
	Create(
		ctx context.Context,
		name string,
	) (*domain.Group, common.SubReturnCode)
	Get(
		ctx context.Context,
		groupId string,
	) (*domain.Group, common.SubReturnCode)
	Update(
		ctx context.Context,
		request domain.Group,
	) common.SubReturnCode
	Delete(ctx context.Context, groupId string) common.SubReturnCode
}

type groupService struct {
	groupRepository repositories2.GroupRepository
	roomRepository  repositories2.RoomRepository
	fileService     FileService
}

func (g *groupService) Delete(ctx context.Context, groupId string) common.SubReturnCode {
	logger := denny.GetLogger(ctx)

	err := g.groupRepository.Delete(ctx, groupId)
	if err != nil {
		logger.WithError(err).Errorln("delete group fail: ", err)
		return common.SystemError
	}

	return common.OK
}

func (g *groupService) Update(
	ctx context.Context,
	request domain.Group,
) common.SubReturnCode {
	var (
		err    = error(nil)
		logger = denny.GetLogger(ctx)
	)
	group, err := g.groupRepository.Get(ctx, request.Id)
	if err != nil {
		logger.WithError(err).Errorln("get group in repository fail: ", err)
		return common.SystemError
	}
	err = g.groupRepository.Update(ctx, models.Group{
		GroupId: request.Id,
		Name:    request.Name,
		FileId:  request.ImageId,
	})
	if err != nil {
		logger.WithError(err).Errorln("update group in repository fail: ", err)
		return common.SystemError
	}
	if request.ImageId != 0 {
		errCode := g.fileService.Delete(ctx, group.FileId)
		if errCode != common.OK {
			logger.Errorln("delete file fail")
			return errCode
		}
	}
	return common.OK
}

func (g *groupService) Get(
	ctx context.Context,
	groupId string,
) (*domain.Group, common.SubReturnCode) {
	var (
		group   = &models.Group{}
		err     = error(nil)
		resp    = &domain.Group{}
		logger  = denny.GetLogger(ctx)
		userIds = make([]string, 0)
	)
	group, err = g.groupRepository.Get(ctx, groupId)
	if err != nil {
		logger.WithError(err).Errorf("find group fail: %s", err)
		return resp, common.SystemError
	}
	rooms, err := g.roomRepository.Get(ctx, domain.Room{
		GroupId: groupId,
	})
	if err != nil {
		logger.WithError(err).Errorf("get room fail: %s", err)
		return resp, common.SystemError
	}
	for _, room := range rooms {
		userIds = append(userIds, room.UserId)
	}
	file, errCode := g.fileService.Get(ctx, group.FileId)
	if errCode != common.OK {
		logger.WithError(err).Errorf("get room fail: %s", err)
		return resp, common.SystemError
	}
	resp = &domain.Group{
		Id:       groupId,
		Name:     group.Name,
		ImageUrl: file.SecureURL,
		UserId:   userIds,
	}

	return resp, common.OK
}

func (g *groupService) Create(
	ctx context.Context,
	name string,
) (*domain.Group, common.SubReturnCode) {
	var (
		resp       = &domain.Group{}
		logger     = denny.GetLogger(ctx)
		groupId, _ = uuid.NewUUID()
	)
	acc := models.Group{
		GroupId: groupId.String(),
		Name:    name,
		FileId:  1,
	}

	err := g.groupRepository.Create(ctx, acc)
	if err != nil {
		logger.Errorln("Create Group service err: ", err)
		return resp, common.SystemError
	}
	resp = &domain.Group{
		Id:   acc.GroupId,
		Name: acc.Name,
	}
	return resp, common.OK
}

func NewGroupService(
	groupRepository repositories2.GroupRepository,
	roomRepository repositories2.RoomRepository,
	fileService FileService,
) GroupService {
	return &groupService{
		groupRepository: groupRepository,
		roomRepository:  roomRepository,
		fileService:     fileService,
	}
}
