package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/models"
	"chat-app/pkg/client/cloudinary"
	"chat-app/pkg/repositories"
	"context"
	"github.com/whatvn/denny"
	"time"
)

type FileService interface {
	Create(ctx context.Context, in domain.UploadIn) (*domain.File, common.SubReturnCode)
	Delete(ctx context.Context, fileId uint) common.SubReturnCode
	Get(ctx context.Context, fileId uint) (*domain.File, common.SubReturnCode)
}

func NewFileService(
	cldClient cloudinary.CloudinaryAPI,
	fileRepository repositories.FileRepository,
) FileService {
	return &fileService{
		cld:            cldClient,
		fileRepository: fileRepository,
	}
}

type fileService struct {
	cld            cloudinary.CloudinaryAPI
	fileRepository repositories.FileRepository
}

func (f *fileService) Get(ctx context.Context, fileId uint) (*domain.File, common.SubReturnCode) {
	logger := denny.GetLogger(ctx)
	resp, err := f.fileRepository.Get(ctx, fileId)
	if err != nil {
		logger.WithError(err).Errorln("get from repository fail: ", err)
		return nil, common.SystemError
	}
	return &domain.File{
		Id:               resp.ID,
		SecureURL:        resp.SecureURL,
		OriginalFilename: resp.OriginalFilename,
	}, common.OK
}

func (f *fileService) Delete(
	ctx context.Context,
	fileId uint,
) common.SubReturnCode {
	logger := denny.GetLogger(ctx)
	err := f.fileRepository.Delete(ctx, fileId)
	if err != nil {
		logger.WithError(err).Errorln("delete file in repository fail: ", err)
		return common.SystemError
	}
	return common.OK
}

func (f *fileService) Create(
	ctx context.Context,
	in domain.UploadIn,
) (*domain.File, common.SubReturnCode) {
	logger := denny.GetLogger(ctx)
	ctxTime, cancel := context.WithTimeout(ctx, 60*time.Second)
	resp, err := f.cld.Upload(ctxTime, in)
	cancel()
	if err != nil {
		logger.WithError(err).Errorln("Create fail err: ", err)
		return nil, common.SystemError
	}
	file := models.File{
		AssetID:          resp.AssetID,
		PublicID:         resp.PublicID,
		AssetFolder:      resp.AssetFolder,
		DisplayName:      resp.DisplayName,
		URL:              resp.URL,
		SecureURL:        resp.SecureURL,
		OriginalFilename: resp.OriginalFilename,
	}
	err = f.fileRepository.Create(ctx, &file)
	if err != nil {
		logger.WithError(err).Errorln("Create fail err: ", err)
		return nil, common.SystemError
	}
	return &domain.File{
		Id:               file.ID,
		SecureURL:        resp.SecureURL,
		OriginalFilename: resp.OriginalFilename,
	}, common.OK
}
