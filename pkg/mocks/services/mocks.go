package services

import (
	"hipeople_task/pkg/domain"
)

type ImageServiceMock struct {
	UploadImageMock func(img *domain.ImageFile) (string, *domain.Error)
	GetImageMock func(id string) ([]byte, *domain.Error)
}


func (i ImageServiceMock) UploadImage(img *domain.ImageFile) (string, *domain.Error) {
	return i.UploadImageMock(img)
}

func (i ImageServiceMock) GetImage(id string) ([]byte, *domain.Error) {
	return i.GetImageMock(id)
}