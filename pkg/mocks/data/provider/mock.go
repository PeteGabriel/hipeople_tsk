package provider

import (
	"io"
)

type ImageProviderMock struct {
	SaveImageMock func(fName string, img io.Reader) (string, error)
	GetImageMock func(imgId string) ([]byte, error)
}


func (m ImageProviderMock) GetImage(imgId string) ([]byte, error) {
	return m.GetImageMock(imgId)
}

func (m ImageProviderMock) SaveImage(fName string, img io.Reader) (string, error) {
	return m.SaveImageMock(fName, img)
}