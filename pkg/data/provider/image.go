package provider

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hipeople_task/pkg/models"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type IImageDataProvider interface {
	SaveImage(img *models.ImageFile) string
	GetImage(imgId string) string
}

type ImageDataProvider struct {
	dataSource map[string]string //[imageId] -> image name
}

func New() IImageDataProvider {
	return &ImageDataProvider{
		dataSource: make(map[string]string),
	}
}

//TODO return errors
func (idp ImageDataProvider) SaveImage(img *models.ImageFile) string{
	//Step 1 - rename file
	imgName := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), img.Header.Filename)
	fP := "./static/images/" + imgName
	//TODO creation of directories must be checked/created at the beginning of program.
	f, err := os.OpenFile(fP, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer f.Close()
	io.Copy(f, img.Content)

	//Step 2 - Create ImageID (hash)
	hsr := sha256.New()
	if _, err := io.Copy(hsr, img.Content); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hsr.Sum(nil)[:32] //amount of bytes produces by sha256
	imageId := hex.EncodeToString(hashInBytes)

	idp.dataSource[imageId] = imgName

	return imageId
}

func (idp ImageDataProvider) GetImage(imgId string) string{
	imgName := idp.dataSource[imgId]

	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile("./static/images/" + imgName)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}