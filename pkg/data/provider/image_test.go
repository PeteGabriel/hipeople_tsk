package provider

import (
	"hipeople_task/pkg/domain"
	"mime/multipart"
	"testing"
)

func TestSaveImage(t *testing.T){

	fh := &multipart.FileHeader{
		Size: 1024*1024,
		Filename: "testing_file.png",
	}

	img := &domain.ImageFile{
		FileName: fh.Filename,
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