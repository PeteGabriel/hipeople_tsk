package provider

import (
	"hipeople_task/pkg/domain"
	mock "hipeople_task/pkg/mocks/data/provider"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveImage(t *testing.T){

	fh := &multipart.FileHeader{
		Size: 1024*1024,
		Filename: "testing_file.png",
	}

	img := &domain.ImageFile{
		FileName: fh.Filename,
		Content: &mock.File{},
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	prov := &ImageDataProvider{
		dataSource:      make(map[string]string),
		storageLocation: filepath.Join(wd, "../../../", relativePath),
	}

	imgId, err := prov.SaveImage(img)

	if err != nil {
		t.Errorf("error not nil: %s", err.Error())
		return
	}

	if len(imgId) == 0 {
		t.Errorf("unexpected image id: %s", imgId)
		return
	}

	//assert internal mapping
	if len(prov.dataSource) <= 0 {
		t.Errorf("unexpected mapping len: %d", len(prov.dataSource))
		return
	}
}