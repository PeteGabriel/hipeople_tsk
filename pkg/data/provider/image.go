package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hipeople_task/pkg/domain"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	relativePath     = "/static/images/"
	ImageNotFoundErr = "image not found in storage"
)

//IImageDataProvider represents the contract this provider gives to all
//entities that make use of it.
type IImageDataProvider interface {
	SaveImage(img *domain.ImageFile) (string, error)
	GetImage(imgId string) ([]byte, error)
}

type ImageDataProvider struct {
	dataSource map[string]string //[imageId] -> image name
}

//NewImageProvider creates a new instance of IImageDataProvider
func NewImageProvider() IImageDataProvider {
	return &ImageDataProvider{
		dataSource: make(map[string]string),
	}
}

//SaveImage saves the image file given as parameter for retrieval later on.
//Renames the file to allow duplicates to be uploaded at any given time.
//It returns an identifier that can be used to retrieve the image.
func (idp ImageDataProvider) SaveImage(img *domain.ImageFile) (string, error) {
	//rename file and store it
	img.FileName = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), img.FileName)
	fPath, err := buildFilePath(img.FileName)
	if err != nil {
		return "", err
	}

	if _, err := copyImage(fPath, img.Content); err != nil {
		return "", err
	}
	//create an ImageID
	imageId, err := createImageId(img)
	if err != nil {
		return "", err
	}

	//map id to image file path
	idp.dataSource[imageId] = fPath

	return imageId, nil
}

func copyImage(fn string, content multipart.File) (bool, error) {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666) // read&write mode
	if err != nil {
		log.Println("error saving image in data provider", err)
		return false, err
	}
	defer f.Close()
	if _, err = io.Copy(f, content); err != nil {
		log.Println("error copying content in data provider", err)
		return false, err
	}
	return true, nil
}

func buildFilePath(filename string) (string, error){
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return filepath.Join(wd, relativePath, filename), nil
}

func createImageId(img *domain.ImageFile) (string, error) {
	hashFn := sha256.New()
	if _, err := io.Copy(hashFn, img.Content); err != nil {
		log.Println("error hashing image", err)
		return "", err
	}
	if _, err := hashFn.Write([]byte(img.FileName)); err != nil { //add the file name to add more entropy
		log.Println("error hashing image", err)
		return "", err
	}
	hashInBytes := hashFn.Sum(nil)[:32] //32 is the amount of bytes produces by sha256
	return hex.EncodeToString(hashInBytes), nil
}

//GetImage by given image ID. Check if given ID is mapped to any file name.
//If so, convert the contents to base64 and return it.
//If not, an error is returned saying that image could not be found.
func (idp ImageDataProvider) GetImage(imgId string) ([]byte, error) {
	fPath := idp.dataSource[imgId]

	if fPath == "" {
		err := fmt.Errorf("%s", ImageNotFoundErr)
		log.Println(err)
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return bytes, nil
}
