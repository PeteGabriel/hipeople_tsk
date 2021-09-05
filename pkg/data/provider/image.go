package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	SaveImage(fName string, img io.Reader) (string, error)
	GetImage(imgId string) ([]byte, error)
}

type ImageDataProvider struct {
	storageLocation string
	dataSource      map[string]string //[imageId] -> image name
}

//NewImageProvider creates a new instance of IImageDataProvider
func NewImageProvider() IImageDataProvider {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err) //can progress if this is not possible
	}
	return &ImageDataProvider{
		dataSource:      make(map[string]string),
		storageLocation: filepath.Join(wd, relativePath),
	}
}

//SaveImage saves the image file given as parameter for retrieval later on.
//Renames the file to allow duplicates to be uploaded at any given time.
//It returns an identifier that can be used to retrieve the image.
func (idp ImageDataProvider) SaveImage(fName string, img io.Reader) (string, error) {
	//rename file and store it
	fPath := idp.buildFilePath(fmt.Sprintf("%d_%s", time.Now().UnixMilli(), fName))

	if _, err := copyImage(fPath, img); err != nil {
		return "", err
	}
	//create an ImageID
	imageId, err := createImageId(fName, img)
	if err != nil {
		return "", err
	}

	//map id to image file path
	idp.dataSource[imageId] = fPath

	return imageId, nil
}

func copyImage(fn string, content io.Reader) (bool, error) {
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

func (idp ImageDataProvider) buildFilePath(filename string) string{
	return filepath.Join(idp.storageLocation, filename)
}

func createImageId(fName string, img io.Reader) (string, error) {
	hashFn := sha256.New()
	if _, err := io.Copy(hashFn, img); err != nil {
		log.Println("error hashing image", err)
		return "", err
	}
	if _, err := hashFn.Write([]byte(fName)); err != nil { //add the file name to add more entropy
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
