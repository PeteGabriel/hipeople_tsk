package middleware

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/models/responses"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

//ValidateContentType validates the file content type to check if user is trying to upload
//something else than an image.
func ValidateContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))
			return
		}
		defer file.Close()

		//check for png/jpeg content-type
		if !isFileExtensionAllowed(handler) {
			log.Println("file type not valid for upload")
			fileExtensionsInvalidErr := responses.ErrProblem{
				Title:    "file received is invalid",
				Detail:   "File type not valid for upload. File received must be an image.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(fileExtensionsInvalidErr.Status)

			res, err := json.Marshal(fileExtensionsInvalidErr)
			if err != nil {
				// TODO handle error
			}
			w.Write(res) //todo handle error
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isFileExtensionAllowed(h *multipart.FileHeader) bool {
	if len(h.Header["Content-Type"]) > 0 {
		contentType := h.Header["Content-Type"][0]
		return strings.HasPrefix(contentType, "image/")
	}else {
		return false
	}
}