package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/client/cloudinary"
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang/glog"
	"time"
)

type MediaService interface {
	Push(in domain.UploadIn) (*uploader.UploadResult, common.SubReturnCode)
}

func NewMediaService(cldClient cloudinary.CloudinaryAPI) MediaService {
	return &mediaServiceImpl{
		cld: cldClient,
	}
}

type mediaServiceImpl struct {
	cld cloudinary.CloudinaryAPI
}

func (m mediaServiceImpl) Push(
	in domain.UploadIn,
) (*uploader.UploadResult, common.SubReturnCode) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	resp, err := m.cld.Upload(ctx, in)
	cancel()
	if err != nil {
		glog.Errorln("Push fail err: ", err)
		return nil, common.SystemError
	}
	return resp, common.OK
}
