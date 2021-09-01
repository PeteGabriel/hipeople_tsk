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

const (
	jsonMediaType        = "application/json"
	problemJsonMediaType = "application/problem+json"
	mediaType = "Content-Type"
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
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))

			w.Header().Set(mediaType, problemJsonMediaType)
			w.WriteHeader(http.StatusBadRequest)
			notReceivedImgErr := responses.ErrProblem{
				Title:    "image not received",
				Detail:   "an image must be provided",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			res, err := json.Marshal(notReceivedImgErr)
			if err != nil {
				// TODO handle error
			}
			w.Write(res) //todo handle error
			return
		}
		defer file.Close()

		imgId, upErr := a.imgService.UploadImage(models.NewImageFile(handler, file))
		if upErr != nil {
			log.Println(fmt.Sprintf("%s - %s", "error uploading received file", upErr.Error.Error()))
			log.Println(fmt.Sprintf("%s - %d - %s", upErr.Name, upErr.Code, upErr.Message))

			var probJson responses.ErrProblem
			if upErr.Name == models.ServerErr {
				w.Header().Set(mediaType, problemJsonMediaType)
				w.WriteHeader(http.StatusInternalServerError)
				probJson = responses.ErrProblem{
					Title:    "image not uploaded",
					Detail:   upErr.Message,
					Status:   http.StatusInternalServerError,
					Instance: r.URL.Path,
				}
			} else {
				//TODO check if theres another possibility than ServerSideError
				w.Header().Set(mediaType, problemJsonMediaType)
				w.WriteHeader(http.StatusBadRequest)
				probJson = responses.ErrProblem{
					Title:    "image not uploaded",
					Detail:   upErr.Message,
					Status:   http.StatusBadRequest,
					Instance: r.URL.Path,
				}
			}

			res, err := json.Marshal(probJson)
			if err != nil {
				// TODO handle error
			}
			w.Write(res) //todo handle error
			return
		}

		w.Header().Set(mediaType, jsonMediaType)
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

			w.Header().Set("Content-Type", problemJsonMediaType)
			w.WriteHeader(http.StatusNotFound)
			ptrnErr := responses.ErrProblem{
				Title:    "image not found",
				Status:   http.StatusNotFound,
				Instance: r.URL.Path,
			}
			res, err := json.Marshal(ptrnErr)
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

			w.Header().Set(mediaType, problemJsonMediaType)
			w.WriteHeader(http.StatusNotFound)
			errProbRes := responses.ErrProblem{
				Title:    "image not found",
				Detail:   err.Message,
				Status:   http.StatusNotFound,
				Instance: r.URL.Path,
			}
			res, err := json.Marshal(errProbRes)
			if err != nil {
				// handle error
			}
			w.Write(res) //todo handle error
			return
		}

		w.Header().Set(mediaType, jsonMediaType)
		w.WriteHeader(http.StatusCreated)
		res, marshalErr := json.Marshal(responses.GetImageResponse{Image: content})
		if marshalErr != nil {
			// handle error
		}
		w.Write(res) //todo handle error

	})
}
