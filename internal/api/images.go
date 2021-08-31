package api

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/models"
	"hipeople_task/pkg/models/responses"
	"log"
	"net/http"
	"regexp"
	"strings"
)

//Upload an image
func (a App) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			//TODO handle better this error
			fmt.Println(err)
			return
		}
		defer file.Close()

		imgId, upErr := a.imgService.UploadImage(models.NewImageFile(handler, file))
		if upErr != nil {
			//TODO Handle this
			panic(upErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		res, err := json.Marshal(responses.UploadResponse{ImageId: imgId})
		if err != nil {
			// handle error
		}
		w.Write(res) //todo handle error
	})
}

//GetImage retrieves an image for a given ID.
func (a App) GetImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//validate info coming from the url path
		if match, _ := regexp.MatchString(`/api/image/[0-9a-z]+$`, r.URL.Path); !match {
			log.Println("Pattern not matched. Image not found.")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			notFoundImgErr := models.Error{
				Message: "image not found",
				Code: http.StatusNotFound,
			}
			res, err := json.Marshal(notFoundImgErr)
			if err != nil {
				// TODO handle error
			}
			w.Write(res) //todo handle error
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		imgId := parts[3] //grab the imageID

		content, err := a.imgService.GetImage(imgId)
		if err != nil {
			if err.Name == models.ServerErr {
				log.Println(fmt.Sprintf("%s - %d - %s", err.Name, err.Code, err.Message))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte{})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(err)
			if err != nil {
				// handle error
			}
			w.Write(res) //todo handle error
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		res, marshalErr := json.Marshal(responses.GetImageResponse{Image: content})
		if marshalErr != nil {
			// handle error
		}
		w.Write(res) //todo handle error

	})
}