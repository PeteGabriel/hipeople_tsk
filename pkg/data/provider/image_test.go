package provider

import (
	"fmt"
	"hipeople_task/pkg/domain"
	"log"
	"mime/multipart"
	"os"
	"testing"
)

func TestSaveImage(t *testing.T){

	fh := &multipart.FileHeader{
		Size: 1024*1024,
		Filename: "testing_file.png",
	}

	img := &domain.ImageFile{
		Header: fh,
		Content: nil,
	}
	prov := NewImageProvider()

	imgId, err := prov.SaveImage(img)

	if err != nil {
		t.Errorf("error not nil: %s", err.Error())
		return
	}

	if len(imgId) == 0 {
		t.Errorf("unexpected image id: %s", imgId)
		return
	}
}