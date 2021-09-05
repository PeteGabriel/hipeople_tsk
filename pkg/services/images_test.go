package services

import (
	"bytes"
	"errors"
	"hipeople_task/pkg/data/provider"
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


func TestImageService_GetImage_ImageNotFound(t *testing.T) {
	providerMock := &mock.ImageProviderMock{}
	providerMock.GetImageMock = func(imgId string) ([]byte, error) {
		return []byte{}, errors.New(provider.ImageNotFoundErr)
	}

	svc := ImageService{provider: providerMock}

	content, err := svc.GetImage("some_good_id")
	if err == nil {
		t.Error("error should not be nil")
		return
	}

	if err.Code != 404 {
		t.Errorf("expected code - 404. Code received: %d", err.Code)
		return
	}

	if err.Error.Error() != provider.ImageNotFoundErr {
		t.Errorf("unexpected internal error message: ´%s´", err.Error.Error())
		return
	}

	if len(content) > 0 {
		t.Error("content should be empty")
		return
	}

}

func TestImageService_GetImage_ServerError(t *testing.T) {
	providerMock := &mock.ImageProviderMock{}
	providerMock.GetImageMock = func(imgId string) ([]byte, error) {
		return []byte{}, errors.New("internal error message")
	}

	svc := ImageService{provider: providerMock}

	content, err := svc.GetImage("some_good_id")

	if err == nil {
		t.Error("error should not be nil")
		return
	}

	if err.Code != 500 {
		t.Errorf("expected code - 500. Code received: %d", err.Code)
		return
	}

	if err.Error.Error() != "internal error message" {
		t.Errorf("unexpected internal error message: ´%s´", err.Error.Error())
		return
	}

	if len(content) > 0 {
		t.Error("content should be empty")
		return
	}
}