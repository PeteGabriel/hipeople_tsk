package api

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/models"
	"hipeople_task/pkg/models/responses"
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

		//check info from url
		if match, _ := regexp.MatchString(`/api/image/[0-9a-z]+$`, r.URL.Path); !match {
			//TODO handle not match flow with JSON error
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("image not found"))
			if err != nil {
				//TODO add to logs
				return
			}
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		imgId := parts[3]

		content, _ := a.imgService.GetImage(imgId)
		w.Write([]byte(content))

	})
}

