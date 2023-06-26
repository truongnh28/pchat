package controller

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
	"github.com/whatvn/denny"
	"net/http"
)

type group struct {
	userService  service.UserService
	groupService service.GroupService
	roomService  service.RoomService
	fileService  service.FileService
}

func (g *group) Create(
	ctx context.Context,
	request *chat_app.CreateGroupRequest,
) (resp *chat_app.CreateGroupResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
		groupResp = &domain.Group{}
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.CreateGroupResponse)
	if request.GetGroupName() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("group_name invalid")
		return
	}
	groupResp, errCode = g.groupService.Create(ctx, request.GetGroupName())
	if errCode != common.OK {
		logger.Errorln("create group fail")
		return
	}
	resp.Info = &chat_app.GroupInfo{
		GroupId:   groupResp.Id,
		GroupName: groupResp.Name,
	}
	return
}

func (g *group) Get(
	ctx context.Context,
	request *chat_app.GroupRequest,
) (resp *chat_app.GetGroupResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
		groupResp = &domain.Group{}
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.GetGroupResponse)
	if request.GetGroupId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("group_id invalid")
		return
	}
	groupResp, errCode = g.groupService.Get(ctx, request.GetGroupId())
	if errCode != common.OK {
		logger.Errorln("get group fail")
		return
	}
	resp.Info = &chat_app.GroupInfo{
		GroupId:   groupResp.Id,
		GroupName: groupResp.Name,
		AvatarUrl: groupResp.ImageUrl,
		UserId:    groupResp.UserId,
	}
	return
}

func (g *group) Update(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.UpdateGroupResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
		ok        = false
		httpCtx   *denny.Context
		uploadRes *domain.File
		imageId   = uint(0)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.UpdateGroupResponse)
	httpCtx, ok = ctx.(*denny.Context)
	if !ok {
		errCode = common.SystemError
		logger.WithError(common.GetHttpCtxFail)
		return
	}
	groupId := httpCtx.Request.FormValue("group_id")
	if groupId == "" {
		errCode = common.InvalidRequest
		logger.WithError(common.FiledInvalid)
		return
	}
	err = httpCtx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		errCode = common.InvalidRequest
		logger.WithError(common.ParseDataFail)
		return
	}
	groupName := httpCtx.Request.FormValue("group_name")
	file, fileHeader, err := httpCtx.Request.FormFile("group_image")
	if err == http.ErrMissingFile && groupName == "" {
		errCode = common.InvalidRequest
		logger.Errorln("Get file from request err: ", err)
		return
	}

	if err == nil {
		uploadRes, errCode = g.fileService.Create(ctx, domain.UploadIn{
			FileName: fileHeader.Filename,
			FileData: file,
		})
		imageId = uploadRes.GetId()
	}

	errCode = g.groupService.Update(ctx, domain.Group{
		Id:      groupId,
		Name:    groupName,
		ImageId: imageId,
	})

	if errCode != common.OK {
		logger.Errorln("update group fail")
		return
	}
	return
}

func (g *group) JoinGroup(
	ctx context.Context,
	request *chat_app.JoinGroupRequest,
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
	if request.GetGroupId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("group_id invalid")
		return
	}
	errCode = g.roomService.Create(ctx, domain.Room{
		GroupId: request.GetGroupId(),
		UserId:  userId,
	})
	if errCode != common.OK {
		logger.Errorln("join group fail")
		return
	}
	return
}

func (g *group) Delete(
	ctx context.Context,
	request *chat_app.GroupRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)
	if request.GetGroupId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("group_id invalid")
		return
	}
	if request.GetGroupId() == "" {
		errCode = common.InvalidRequest
		logger.Errorln("group id is not empty")
		return
	}
	errCode = g.groupService.Delete(ctx, request.GetGroupId())
	if errCode != common.OK {
		logger.Errorln("delete group fail")
		return
	}

	return
}

func NewGroup(
	userService service.UserService,
	groupService service.GroupService,
	roomService service.RoomService,
	fileService service.FileService,
) chat_app.GroupServer {
	return &group{
		userService:  userService,
		groupService: groupService,
		roomService:  roomService,
		fileService:  fileService,
	}
}
