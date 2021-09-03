package provider

import "hipeople_task/pkg/domain"

type ImageProviderMock struct {
	SaveImageMock func(img *domain.ImageFile) (string, error)
	GetImageMock func(imgId string) ([]byte, error)
}


func (m ImageProviderMock) GetImage(imgId string) ([]byte, error) {
	return m.GetImageMock(imgId)
}

func (m ImageProviderMock) SaveImage(img *domain.ImageFile) (string, error) {
	return m.SaveImageMock(img)
}