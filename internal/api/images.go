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
	imageMediaType       = "image/*"
	fileSizeLimit        = 1024 * 1024
	getImgRouteRegex     = `/api/image/[0-9a-z]+$`
)

//Upload an image
func (a App) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, fileSizeLimit) //limit the size of the request body
		//ask for multipart to be parsed against a certain limit size.
		if err := r.ParseMultipartForm(fileSizeLimit); err != nil {
			log.Println("error parsing multipart form: ", err)
			probJson := responses.ErrProblem{
				Title:    "image not uploaded",
				Detail:   "Image not uploaded. Check file size (1MB max) or upload form.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			writeProblemJson(w, probJson)
			return
		}

		file, handler, err := r.FormFile("upload_file")
		if err != nil {
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))
			return
		}
		defer file.Close()

		imgId, upErr := a.imgService.UploadImage(domain.NewImageFile(handler.Filename, file))
		if upErr != nil {
			log.Println(fmt.Sprintf("%s - %s", "error uploading received file", upErr.Error.Error()))
			log.Println(fmt.Sprintf("%s - %d - %s", upErr.Name, upErr.Code, upErr.Message))

			w.Header().Set(mediaType, "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
			return
		}

		w.Header().Set(mediaType, jsonMediaType)
		w.WriteHeader(http.StatusCreated)
		res, err := json.Marshal(responses.UploadResponse{ImageId: imgId})
		if err != nil {
			log.Println("error marshaling uploaded image response", err)
		}
		w.Write(res)
	}
}

//GetImage retrieves an image for a given ID.
func (a App) GetImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//validate info coming from the url path
		if match, _ := regexp.MatchString(getImgRouteRegex, r.URL.Path); !match {
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
				w.Write([]byte{})
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

		// Send image as response
		w.Header().Set(mediaType, imageMediaType)
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}

func writeProblemJson(w http.ResponseWriter, prob responses.ErrProblem) {
	w.Header().Set(mediaType, problemJsonMediaType)
	w.WriteHeader(prob.Status)

	res, err := json.Marshal(prob)
	if err != nil {
		log.Println("error marshaling uploaded image response", err)
	}
	w.Write(res)
}


