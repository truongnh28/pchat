package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/models"
	"chat-app/pkg/repositories"
	chat_app "chat-app/proto/chat-app"
	"context"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/whatvn/denny"
)

type GroupService interface {
	Create(
		ctx context.Context,
		req *chat_app.CreateGroupRequest,
	) (*chat_app.CreateGroupResponse, common.SubReturnCode)
	Get(
		ctx context.Context,
		req *chat_app.GetGroupRequest,
	) (*chat_app.GetGroupResponse, common.SubReturnCode)
}

type groupService struct {
	groupRepository repositories.GroupRepository
	roomRepository  repositories.RoomRepository
}

func (a *groupService) Get(
	ctx context.Context,
	req *chat_app.GetGroupRequest,
) (*chat_app.GetGroupResponse, common.SubReturnCode) {
	var (
		acc     = &models.Group{}
		err     = error(nil)
		resp    = &chat_app.GetGroupResponse{}
		logger  = denny.GetLogger(ctx)
		userIds = make([]string, 0)
	)
	if req.GetGroupId() == "" {
		return resp, common.SystemError
	}
	acc, err = a.groupRepository.Get(ctx, req.GetGroupId())
	if err != nil {
		glog.Errorf("find group fail: %s", err)
		logger.WithError(err)
		return resp, common.SystemError
	}
	rooms, err := a.roomRepository.Get(ctx, domain.Room{
		GroupId: req.GetGroupId(),
	})
	if err != nil {
		glog.Errorf("get room fail: %s", err)
		logger.WithError(err)
		return resp, common.SystemError
	}
	for _, room := range rooms {
		userIds = append(userIds, room.UserId)
	}
	resp.Info = &chat_app.GroupInfo{
		GroupId:   req.GroupId,
		GroupName: acc.Name,
		AvatarUrl: acc.AvatarUrl,
		UserId:    userIds,
	}
	return resp, common.OK
}

func (a *groupService) Create(
	ctx context.Context,
	req *chat_app.CreateGroupRequest,
) (*chat_app.CreateGroupResponse, common.SubReturnCode) {
	var (
		resp       = &chat_app.CreateGroupResponse{}
		logger     = denny.GetLogger(ctx)
		groupId, _ = uuid.NewUUID()
	)
	acc := models.Group{
		GroupId:   groupId.String(),
		Name:      req.GetGroupName(),
		AvatarUrl: req.GetAvatarUrl(),
	}

	err := a.groupRepository.Create(ctx, acc)
	if err != nil {
		logger.Errorln("Create Group service err: ", err)
		return resp, common.SystemError
	}
	resp.Info = &chat_app.GroupInfo{
		GroupId:   groupId.String(),
		GroupName: req.GetGroupName(),
		AvatarUrl: req.GetAvatarUrl(),
	}
	return resp, common.OK
}

func NewGroupService(
	groupRepository repositories.GroupRepository,
	roomRepository repositories.RoomRepository,
) GroupService {
	return &groupService{
		groupRepository: groupRepository,
		roomRepository:  roomRepository,
	}
}
