package services

import (
	"hipeople_task/pkg/data/provider"
	"hipeople_task/pkg/models"
)

//IImageService represents the contract this service gives to all
//entities that make use of it.
type IImageService interface {
	UploadImage(img *models.ImageFile) (string, *models.Error)
	GetImage(id string) (string, *models.Error)
}

type ImageService struct {
	provider provider.IImageDataProvider
}

//New creates a new instance of IImageService
func New() IImageService {
	return &ImageService{
		provider: provider.New(),
	}
}

//UploadImage uploads an image described by the image file given as parameter.
func (i ImageService) UploadImage(img *models.ImageFile) (string, *models.Error) {
	imgId := i.provider.SaveImage(img)
	return imgId, nil
}

//GetImage by image id.
func (i ImageService) GetImage(id string) (string, *models.Error) {
	img, err := i.provider.GetImage(id)
	if err != nil {
		//todo handle this
	}

	return img, nil
}