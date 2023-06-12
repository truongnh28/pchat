package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/client"
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang/glog"
	"time"
)

type MediaService interface {
	Upload(in domain.UploadIn) (*uploader.UploadResult, common.SubReturnCode)
}

func NewMediaService(cldClient client.CloudinaryAPI) MediaService {
	return &mediaServiceImpl{
		cld: cldClient,
	}
}

type mediaServiceImpl struct {
	cld client.CloudinaryAPI
}

func (m mediaServiceImpl) Upload(
	in domain.UploadIn,
) (*uploader.UploadResult, common.SubReturnCode) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	resp, err := m.cld.Upload(ctx, in)
	cancel()
	if err != nil {
		glog.Errorln("Upload fail err: ", err)
		return nil, common.SystemError
	}
	return resp, common.OK
}
