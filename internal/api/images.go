package api

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/domain"
	"hipeople_task/pkg/models/responses"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	jsonMediaType        = "application/json"
	problemJsonMediaType = "application/problem+json"
	mediaType            = "Content-Type"
	fileSizeLimit        = 1024 * 1024
)

//Upload an image
func (a App) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, fileSizeLimit) //limit the size of the request body
		if err := r.ParseMultipartForm(fileSizeLimit); err != nil {
			log.Println("file bigger than 1MB: ", err)
			probJson := responses.ErrProblem{
				Title:    "image not uploaded",
				Detail:   "Image not uploaded. File is bigger than 1MB.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			writeProblemJson(w, probJson)
			return
		}
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))
			notReceivedImgErr := responses.ErrProblem{
				Title:    "image not received",
				Detail:   "An image must be provided.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			writeProblemJson(w, notReceivedImgErr)
			return
		}
		defer file.Close()

		imgId, upErr := a.imgService.UploadImage(domain.NewImageFile(handler, file))
		if upErr != nil {
			log.Println(fmt.Sprintf("%s - %s", "error uploading received file", upErr.Error.Error()))
			log.Println(fmt.Sprintf("%s - %d - %s", upErr.Name, upErr.Code, upErr.Message))

			var probJson responses.ErrProblem
			if upErr.Name == domain.ServerErr {
				probJson = responses.ErrProblem{
					Title:    "image not uploaded",
					Detail:   upErr.Message,
					Status:   http.StatusInternalServerError,
					Instance: r.URL.Path,
				}
			} else {
				//TODO check if theres another possibility than ServerSideError
				probJson = responses.ErrProblem{
					Title:    "image not uploaded",
					Detail:   upErr.Message,
					Status:   http.StatusBadRequest,
					Instance: r.URL.Path,
				}
			}

			writeProblemJson(w, probJson)
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
			ptrnErr := responses.ErrProblem{
				Title:    "image not found",
				Status:   http.StatusNotFound,
				Instance: r.URL.Path,
			}
			writeProblemJson(w, ptrnErr)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		imgId := parts[3] //grab the imageID

		content, err := a.imgService.GetImage(imgId)
		if err != nil {
			if err.Name == domain.ServerErr {
				log.Println(fmt.Sprintf("%s - %d - %s", err.Name, err.Code, err.Message))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte{}) //TODO recheck this
				return
			}

			errProbRes := responses.ErrProblem{
				Title:    "image not found",
				Detail:   err.Message,
				Status:   http.StatusNotFound,
				Instance: r.URL.Path,
			}
			writeProblemJson(w, errProbRes)
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

func writeProblemJson(w http.ResponseWriter, prob responses.ErrProblem) {
	w.Header().Set(mediaType, problemJsonMediaType)
	w.WriteHeader(prob.Status)

	res, err := json.Marshal(prob)
	if err != nil {
		// TODO handle error
	}
	w.Write(res) //todo handle error
}
