package services

import (
	"fmt"
	"hipeople_task/pkg/data/provider"
	"hipeople_task/pkg/domain"
)

//IImageService represents the contract this service gives to all
//entities that make use of it.
type IImageService interface {
	UploadImage(img *domain.ImageFile) (string, *domain.Error)
	GetImage(id string) (string, *domain.Error)
}

type ImageService struct {
	provider provider.IImageDataProvider
}

//New creates a new instance of IImageService
func New() IImageService {
	return &ImageService{
		provider: provider.NewImageProvider(),
	}
}

//UploadImage uploads an image described by the image file given as parameter.
func (i ImageService) UploadImage(img *domain.ImageFile) (string, *domain.Error) {
	imgId, err := i.provider.SaveImage(img)
	if err != nil {
		return "", &domain.Error{
			Error:   err,
			Message: fmt.Sprintf("%s %s", "error saving image", img.Header.Filename),
			Code:    500,
			Name:    domain.ServerErr,
		}
	}
	return imgId, nil
}

//GetImage by image id.
func (i ImageService) GetImage(id string) (string, *domain.Error) {
	img, err := i.provider.GetImage(id)
	if err != nil {

		if err.Error() == provider.ImageNotFoundErr {
			return "", &domain.Error{
				Error:   err,
				Message: fmt.Sprintf("%s %s - %s", "error getting image with id", id, err.Error()),
				Code:    404,
			}
		}

		return "", &domain.Error{
			Error:   err,
			Message: fmt.Sprintf("%s %s", "error getting image with id", id),
			Code:    500,
			Name:    domain.ServerErr,
		}
	}

	return img, nil
}
