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

const RelativePath = "./static/images/"

//IImageDataProvider represents the contract this provider gives to all
//entities that make use of it.
type IImageDataProvider interface {
	SaveImage(img *models.ImageFile) string
	GetImage(imgId string) string
}

type ImageDataProvider struct {
	dataSource map[string]string //[imageId] -> image name
}

//New creates a new instance of IImageDataProvider
func New() IImageDataProvider {
	return &ImageDataProvider{
		dataSource: make(map[string]string),
	}
}


//SaveImage saves the image file given as parameter for retrieval later on.
//It renames the file to allow duplicates to be uploaded at any given time.
//It returns an identifier that can be used to retrieve the image.
func (idp ImageDataProvider) SaveImage(img *models.ImageFile) string {
	//TODO return errors
	//Step 1 - rename file
	img.Header.Filename = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), img.Header.Filename)
	fP := RelativePath + img.Header.Filename
	//TODO creation of directories must be checked/created at the beginning of program.
	f, err := os.OpenFile(fP, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer f.Close()
	io.Copy(f, img.Content)

	//Step 2 - Create ImageID
	imageId := createImageId(img)
	//map id to image
	idp.dataSource[imageId] = img.Header.Filename

	return imageId
}

func createImageId(img *models.ImageFile) string {
	hsr := sha256.New()
	if _, err := io.Copy(hsr, img.Content); err != nil {
		log.Fatal(err)
		//todo send error above
	}
	hsr.Write([]byte(img.Header.Filename))//add the file name to add more entropy
	hashInBytes := hsr.Sum(nil)[:32] //amount of bytes produces by sha256
	return hex.EncodeToString(hashInBytes)
}

func (idp ImageDataProvider) GetImage(imgId string) string {
	imgName := idp.dataSource[imgId]

	bytes, err := ioutil.ReadFile("./static/images/" + imgName)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}