package services

import (
	"hipeople_task/pkg/data/provider"
	"hipeople_task/pkg/models"
)

type IImageService interface {
	UploadImage(img *models.ImageFile) (string, *models.Error)
	GetImage(id string) (string, *models.Error)
}

type ImageService struct {
	provider provider.IImageDataProvider
}

func NewImageService() IImageService {
	return &ImageService{
		provider: provider.New(),
	}
}

func (i ImageService) UploadImage(img *models.ImageFile) (string, *models.Error) {
	imgId := i.provider.SaveImage(img)
	return imgId, nil
}


func (i ImageService) GetImage(id string) (string, *models.Error) {
	return i.provider.GetImage(id), nil
	//TODO handle errors and wrapping of them
}