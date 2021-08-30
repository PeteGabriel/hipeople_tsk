package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hipeople_task/pkg/models"
	"io"
	"log"
	"os"
	"time"
)

type IImageService interface {
	UploadImage(img *models.ImageFile) (string, *models.Error)
}

type ImageService struct {

}

func NewImageService() IImageService {
	return &ImageService{}
}

func (i ImageService) UploadImage(img *models.ImageFile) (string, *models.Error) {

	//Step 1 - rename file
	fn := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), img.Header.Filename)
	fP := "./static/images/" + fn
	//TODO creation of directories must be checked/created at the beginning of program.
	f, err := os.OpenFile(fP, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer f.Close()
	io.Copy(f, img.Content)

	//Step 2 - Create ImageID (hash)
	hsr := sha256.New()
	if _, err := io.Copy(hsr, img.Content); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hsr.Sum(nil)[:32] //amount of bytes produces by sha256
	returnSHA1String := hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil
}