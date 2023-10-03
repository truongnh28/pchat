package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/repositories"
	"chat-app/models"
	"context"
	"github.com/whatvn/denny"
)

type FriendService interface {
	AddFriend(ctx context.Context, request domain.Friend) common.SubReturnCode
	GetAllFriendByUserId(ctx context.Context, userId string) ([]domain.User, common.SubReturnCode)
	RemoveFriend(ctx context.Context, request domain.Friend) common.SubReturnCode
}

type friendService struct {
	friendRepository repositories.FriendRepository
	userService      UserService
}

func (f friendService) AddFriend(ctx context.Context, request domain.Friend) common.SubReturnCode {
	var (
		logger = denny.GetLogger(ctx)
	)

	err := f.friendRepository.Create(ctx, models.Friend{
		UserId1: request.UserId1,
		UserId2: request.UserId2,
	})
	if err != nil {
		logger.Errorln("Create friend service err: ", err)
		return common.SystemError
	}
	return common.OK
}

func (f friendService) GetAllFriendByUserId(
	ctx context.Context,
	userId string,
) ([]domain.User, common.SubReturnCode) {
	var (
		friends = make([]models.Friend, 0)
		err     = error(nil)
		resp    = make([]domain.User, 0)
		logger  = denny.GetLogger(ctx)
	)
	friends, err = f.friendRepository.GetAllByUserId(ctx, userId)
	if err != nil {
		logger.WithError(err).Errorf("find friends fail: %s", err)
		return resp, common.SystemError
	}
	for _, i := range friends {
		userIdReq := userId
		if userId == i.UserId1 {
			userIdReq = i.UserId2
		} else {
			userIdReq = i.UserId1
		}
		user, errCode := f.userService.GetByUserId(ctx, userIdReq)
		if errCode != common.OK {
			logger.WithError(err).Errorf("find friends details fail: %s", err)
			//return resp, common.SystemError
			continue
		}
		resp = append(resp, user)
	}
	return resp, common.OK
}

func (f friendService) RemoveFriend(
	ctx context.Context,
	request domain.Friend,
) common.SubReturnCode {
	logger := denny.GetLogger(ctx)

	err := f.friendRepository.Delete(ctx, models.Friend{
		UserId1: request.UserId1,
		UserId2: request.UserId2,
	})
	if err != nil {
		logger.WithError(err).Errorln("delete friend fail: ", err)
		return common.SystemError
	}

	return common.OK
}

func NewFriendService(
	friendRepository repositories.FriendRepository,
	userService UserService,
) FriendService {
	return &friendService{
		friendRepository: friendRepository,
		userService:      userService,
	}
}
