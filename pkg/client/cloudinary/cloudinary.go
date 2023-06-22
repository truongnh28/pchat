package cloudinary

import (
	"chat-app/config"
	"chat-app/internal/domain"
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang/glog"
	"sync"
)

type CloudinaryAPI interface {
	Upload(ctx context.Context, in domain.UploadIn) (*uploader.UploadResult, error)
}
type cloudinaryAPI struct {
	client *cloudinary.Cloudinary
}

var cloudinaryInstance *cloudinaryAPI
var cloudinaryMutex sync.Mutex

func GetCloudinaryAPI(config *config.CloudinaryConfig) CloudinaryAPI {
	if cloudinaryInstance != nil {
		return cloudinaryInstance
	}
	cloudinaryMutex.Lock()
	defer cloudinaryMutex.Unlock()
	var cldUrl = fmt.Sprintf("cloudinary://%s:%s@%s", config.APIKey, config.APISecret, config.Name)
	var cld, err = cloudinary.NewFromURL(cldUrl)
	if err != nil {
		//log.Fatalf("Failed to intialize Cloudinary, %v", err)
		panic(fmt.Errorf("unable to connect to cloudinary: %v", err.Error()))
	}
	cloudinaryInstance = &cloudinaryAPI{
		client: cld,
	}
	return cloudinaryInstance
}

func (c *cloudinaryAPI) Upload(
	ctx context.Context,
	in domain.UploadIn,
) (*uploader.UploadResult, error) {
	resp, err := c.client.Upload.Upload(ctx, in.FileData, uploader.UploadParams{})
	if err != nil {
		glog.Errorln("Push cloudinary fail: ", err)
		return nil, err
	}
	return resp, nil
}
