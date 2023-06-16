package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/client/firebase"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type NotificationService interface {
	Push(in domain.UploadIn) (*uploader.UploadResult, common.SubReturnCode)
}

func NewNotificationService(fb firebase.Firebase) NotificationService {
	return &notificationServiceImpl{
		fb: fb,
	}
}

type notificationServiceImpl struct {
	fb firebase.Firebase
}

func (m notificationServiceImpl) Upload(
	in domain.UploadIn,
) (*uploader.UploadResult, common.SubReturnCode) {

}
