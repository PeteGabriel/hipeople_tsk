package services

import (
	"bytes"
	mock "hipeople_task/pkg/mocks/data/provider"
	"testing"
)

func TestImageService_GetImage(t *testing.T) {

	expectedImageContent := []byte{10,2,23,1,23,11}
	providerMock := &mock.ImageProviderMock{}
	providerMock.GetImageMock = func(imgId string) ([]byte, error) {
		return expectedImageContent, nil
	}

	svc := ImageService{provider: providerMock}

	content, err := svc.GetImage("some_good_id")
	if err != nil {
		t.Errorf("error not nil: %s", err.Error)
		return
	}

	if c := bytes.Compare(expectedImageContent, content); c != 0 {
		t.Errorf("contents are not equal: %d", c)
		return
	}

}
