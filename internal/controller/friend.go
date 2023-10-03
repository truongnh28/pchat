package controller

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
)

type friend struct {
	friendService service.FriendService
}

func (f *friend) AddFriend(
	ctx context.Context,
	request *chat_app.FriendRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)
	if request.GetUserId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("userId is empty")
		return
	}
	errCode = f.friendService.AddFriend(ctx, domain.Friend{
		UserId1: userId,
		UserId2: request.GetUserId(),
	})
	if errCode != common.OK {
		logger.Errorln("add friend fail")
		return
	}
	return
}

func (f *friend) GetAllFriendByUserId(
	ctx context.Context,
	request *chat_app.FriendRequest,
) (resp *chat_app.GetFriendResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
		friends   = make([]domain.User, 0)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.GetFriendResponse)
	if request.GetUserId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("userId is empty")
		return
	}
	friends, errCode = f.friendService.GetAllFriendByUserId(ctx, request.GetUserId())
	if errCode != common.OK {
		logger.Errorln("get all friend fail")
		return
	}
	for _, i := range friends {
		resp.UserInfos = append(resp.UserInfos, &chat_app.UserInfo{
			Username:    i.Username,
			PhoneNumber: i.PhoneNumber,
			Email:       i.Email,
			Status:      i.Status,
			UserId:      i.UserId,
			DateOfBirth: i.DateOfBirth.String(),
			Gender:      i.Gender.String(),
			Url:         i.Url,
		})
	}
	return
}

func (f *friend) RemoveFriend(
	ctx context.Context,
	request *chat_app.FriendRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)
	if request.GetUserId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("userId is empty")
		return
	}
	errCode = f.friendService.RemoveFriend(ctx, domain.Friend{
		UserId1: userId,
		UserId2: request.GetUserId(),
	})
	if errCode != common.OK {
		logger.Errorln("remove friend fail")
		return
	}
	return
}

func NewFriend(
	friendService service.FriendService,
) chat_app.FriendsServer {
	return &friend{
		friendService: friendService,
	}
}
