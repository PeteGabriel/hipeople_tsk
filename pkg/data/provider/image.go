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

const (
	relativePath = "./static/images/"
	ImageNotFoundErr = "image not found"
)
//IImageDataProvider represents the contract this provider gives to all
//entities that make use of it.
type IImageDataProvider interface {
	SaveImage(img *models.ImageFile) (string, error)
	GetImage(imgId string) (string, error)
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
func (idp ImageDataProvider) SaveImage(img *models.ImageFile) (string, error) {
	//TODO return errors
	//Step 1 - rename file and store it
	if _, err := saveImage(img); err != nil {
		log.Println("error saving image", err)
		return "", err
	}
	//Step 2 - Create ImageID
	imageId := createImageId(img)
	//map id to image
	idp.dataSource[imageId] = img.Header.Filename

	return imageId, nil
}

func saveImage(img *models.ImageFile) (bool, error) {
	img.Header.Filename = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), img.Header.Filename)
	fP := relativePath + img.Header.Filename

	f, err := os.OpenFile(fP, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer f.Close()
	io.Copy(f, img.Content)
	return true, nil
}

func createImageId(img *models.ImageFile) string {
	hsr := sha256.New()
	if _, err := io.Copy(hsr, img.Content); err != nil {
		log.Fatal(err)
		//todo send error above
	}
	hsr.Write([]byte(img.Header.Filename)) //add the file name to add more entropy
	hashInBytes := hsr.Sum(nil)[:32]       //amount of bytes produces by sha256
	return hex.EncodeToString(hashInBytes)
}

//GetImage by given image ID. Check if given ID is mapped to any file name.
//If so, convert the contents to base64 and return it.
//If not, an error is returned saying that image could not be found.
func (idp ImageDataProvider) GetImage(imgId string) (string, error) {
	imgName := idp.dataSource[imgId]

	if imgName == "" {
		err := fmt.Errorf("%s", ImageNotFoundErr)
		log.Println(err.Error())
		return "", err
	}

	bytes, err := ioutil.ReadFile(relativePath + imgName)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
