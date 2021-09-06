package middleware

import (
	"encoding/json"
	"fmt"
	"hipeople_task/pkg/models/responses"
	"log"
	"net/http"
)

//CheckReceivedContent check the existence of a file to upload in the request.
func CheckReceivedContent(next http.Handler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("upload_file")
		if err != nil {
			log.Println(fmt.Sprintf("%s - %s", "error receiving file to upload", err.Error()))
			notReceivedImgErr := responses.ErrProblem{
				Title:    "image not received",
				Detail:   "An image must be provided.",
				Status:   http.StatusBadRequest,
				Instance: r.URL.Path,
			}
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(notReceivedImgErr.Status)

			res, err := json.Marshal(notReceivedImgErr)
			if err != nil {
				log.Println("error marshaling response in receiver middleware: ", err)
			}
			w.Write(res)
			return
		}
		defer file.Close()

		if next != nil {
			next.ServeHTTP(w, r)
		}
	}
}
